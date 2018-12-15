package batcher

import (
	"context"
	"fmt"
	"time"

	"github.com/beinan/gql-server/concurrent/cache"
	"github.com/beinan/gql-server/concurrent/future"
)

const MaxBatchSize = 20
const BatchInterval = 3 * time.Millisecond

type ID = string
type Context = context.Context

type Value = interface{}

type Result struct {
	Value Value
	Err   error
}

type GetFunc = func(Context, ID) future.Future
type BatchFunc = func(Context, []ID) []Result

type Batcher struct {
	ctx       Context
	cache     *cache.Cache
	batchFunc BatchFunc
	reqChan   chan req
}

type req struct {
	id       ID
	response chan Result
}

func mkBatcher(ctx Context, batchFunc BatchFunc) *Batcher {
	batcher := &Batcher{
		ctx:       ctx,
		cache:     cache.MkCache(),
		batchFunc: batchFunc,
		reqChan:   make(chan req),
	}
	go batcher.startBatcher()
	return batcher
}

func (b *Batcher) AddToBatch(id ID) future.Future {
	request := req{
		id:       id,
		response: make(chan Result),
	}
	producer := func() (Value, error) {
		b.reqChan <- request
		result := <-request.response
		return result.Value, result.Err
	}
	return future.MakeFuture(producer)
}

func (b *Batcher) startBatcher() {
	var ids []ID
	var responses []chan Result
	reset := func() {
		ids = make([]ID, 0, MaxBatchSize)
		responses = make([]chan Result, 0, MaxBatchSize)
	}
	reset()
	for {
		select {
		case r := <-b.reqChan:
			ids = append(ids, r.id)
			responses = append(responses, r.response)
			if len(ids) == MaxBatchSize {
				go b.callBatchFunc(ids, responses)
				reset()
			}
		case <-time.After(BatchInterval):
			go b.callBatchFunc(ids, responses)
			reset()
		case <-b.ctx.Done():
			//terminate when the client's connection closes,
			//the request is canceled (with HTTP/2),
			//or when the ServeHTTP method returns.
			return
		}
	}
}

//query result for each individual id will be sent back via the channel in responses
func (b *Batcher) callBatchFunc(ids []ID, responses []chan Result) {
	if len(ids) == 0 {
		return //nothing in batch, do nothing
	}
	results := b.batchFunc(b.ctx, ids)
	if len(results) != len(ids) {
		panic(fmt.Sprintf("Unexpected batch query resultsids:%v results:%v",
			ids, results))
	}
	for i, result := range results {
		responses[i] <- result //send back result
	}
}
