package batcher

import (
	"context"
	"sync"

	"github.com/beinan/gql-server/concurrent/future"
)

// Manage all the bathers under one context.
type BatcherManager struct {
	batcherMetas []*BatcherMeta
	mux          sync.Mutex
}

type BatcherContextKey int

type BatcherMeta struct {
	key       BatcherContextKey
	batchFunc BatchFunc
}

func MakeBatcherManager() *BatcherManager {
	return &BatcherManager{
		//make an batcher meta array with 0 length and 20 cap
		batcherMetas: make([]*BatcherMeta, 0, 20),
	}
}

func (b *BatcherManager) mkBatchMeta(batchFunc BatchFunc) *BatcherMeta {
	b.mux.Lock()
	defer b.mux.Unlock()
	meta := &BatcherMeta{
		key:       BatcherContextKey(len(b.batcherMetas)),
		batchFunc: batchFunc,
	}
	b.batcherMetas = append(b.batcherMetas, meta)
	return meta
}

//tho this func isgoroutine safe,
//we still recommmend you register all the batcher in main func/goroutine
func (b *BatcherManager) Register(batchFunc BatchFunc) GetFunc {
	batcherMeta := b.mkBatchMeta(batchFunc)
	return func(ctx Context, id ID) future.Future {
		batcher := batcherMeta.getBatcher(ctx) //get batcher per request
		return batcher.cache.LoadOrElse(id, func() (Value, error) {
			return batcher.AddToBatch(id).Value()
		})
	}
}

//Calling this function in http middleware
func (b *BatcherManager) Attach(ctx Context) Context {
	for _, meta := range b.batcherMetas {
		ctx = context.WithValue(ctx, meta.key, mkBatcher(ctx, meta.batchFunc))
	}
	return ctx
}

func (meta *BatcherMeta) getBatcher(ctx Context) *Batcher {
	return ctx.Value(meta.key).(*Batcher)
}
