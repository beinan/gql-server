package batcher

import (
	"context"
	"sync"

	"github.com/beinan/gql-server/concurrent/future"
)

const MaxBatchSize = 20

type ID = String
type Context = context.Context

type Value = interface{}

type ValueErr struct {
	Value interface{}
	Err   error
}

type Pair struct {
	key   ID
	value Value
}

type req struct {
	id       ID
	response <-chan *Value
}

type Batch struct {
	ids []ID
}

type Cache struct {
	cache map[ID]future.Future //request-scoped cache
	mutex sync.RWMutex
}

func mkCache() *Cache {
	return &Cache{
		cache: make(map[ID]future.Future),
	}
}

type BatcherManager struct {
	batcherCount int
}

//func addIDToBatch(id) *Batch {}
//func (*Batch) getValueByID{} Value

type GetFunc = func(Context, ID) future.Future
type BatchFunc = func(Context, []ID) ([]ValueErr, error)

func (b *BatcherManager) Register(batchFunc BatchFunc) GetFunc {
	batcherIndex := b.batcherCount
	b.batcherCount++
	return func(ctx Context, id ID) future.Future {
		//todo: passed in a context?
		//todo: check if id exists in cache
		return batch.get(id)
	}
}

func (b *BatcherManager) Attach(ctx Context) Context {
	for i := 0; i < b.batcherCount; i++ {
		ctx := context.WithValue(ctx, contextKey(i), mkCache())
	}
	return ctx
}

func (b *BatcherManager) getBatcherCache(ctx Context, index int) *Cache {
	return ctx.Value(contextKey(index)).(*Cache)
}

func contextKey(index int) string {
	return "Batcher_" + i
}
