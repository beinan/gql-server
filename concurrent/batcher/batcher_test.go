package batcher

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

type MockedDAO struct {
	batchCounter int
	idCounter    int
	mutex        sync.Mutex
}

type MockedData struct {
	id   ID
	data string
}

func (dao *MockedDAO) GetByIds(ctx Context, ids []ID) []Result {
	dao.mutex.Lock()
	fmt.Println("calling GetBytIds", ids)
	dao.batchCounter++
	dao.mutex.Unlock()
	time.Sleep(time.Duration(2 * BatchInterval))
	results := make([]Result, len(ids))
	for i, id := range ids {
		dao.mutex.Lock()
		dao.idCounter++
		dao.mutex.Unlock()
		results[i] = Result{
			Value: MockedData{
				id:   id,
				data: "data of " + id,
			},
			Err: nil,
		}
	}
	return results
}

var dao = MockedDAO{}

func TestBatcherMergeSameIds(t *testing.T) {
	manager := MakeBatcherManager()
	dao := MockedDAO{}
	getById := manager.Register(dao.GetByIds)
	ctx := context.Background()
	ctx = manager.Attach(ctx)
	getById(ctx, "1")
	getById(ctx, "2")
	getById(ctx, "2")
	getById(ctx, "1").Value()

	if dao.batchCounter != 1 {
		t.Fatalf("Expected 1 batch sent, but actual is %v", dao.batchCounter)
	}

	if dao.idCounter != 2 {
		t.Fatalf("Expected 2 id sent, but actual is %v", dao.idCounter)
	}
}

func TestBatcherMoreBatches(t *testing.T) {
	manager := MakeBatcherManager()
	dao := MockedDAO{}
	getById := manager.Register(dao.GetByIds)
	ctx := context.Background()
	ctx = manager.Attach(ctx)

	for i := 0; i < MaxBatchSize; i++ {
		getById(ctx, strconv.Itoa(i))
	}
	getById(ctx, "OneExtraID")
	time.Sleep(time.Duration(10 * BatchInterval))

	//should be in two batch
	if dao.batchCounter != 2 {
		t.Fatalf("Expected 2 batch sent, but actual is %v", dao.batchCounter)
	}

	if dao.idCounter != MaxBatchSize+1 {
		t.Fatalf("Expected %v ids sent, but actual is %v", MaxBatchSize+1, dao.idCounter)
	}
}
