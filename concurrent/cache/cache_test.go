package cache

import (
	"sync"
	"testing"
	"time"
)

type MockedDAO struct {
	Counter int
}

type MockedData struct {
	id   ID
	data string
}

func (dao *MockedDAO) GetById(id ID) (MockedData, error) {
	dao.Counter++
	time.Sleep(time.Duration(300 * time.Millisecond))
	return MockedData{id, "test data"}, nil
}

var dao = MockedDAO{}

func TestLoadOrElse(t *testing.T) {
	cache := MkCache()
	producer := func() (Value, error) {
		return dao.GetById("123")
	}
	fv := cache.LoadOrElse("123", producer)
	same_fv, _ := cache.Load("123")
	value, err := fv.Value()
	if data, ok := value.(MockedData); ok && data.id != "123" {
		t.Fatalf("Returned value is incorrect: %v", data)
	}

	if fv != same_fv {
		t.Fatalf("Returned different future value instances for the same key")
	}
	if err != nil {
		t.Fatalf("Returned error is not nil: %v", err)
	}
}

func TestLoadOrElseConcurrent(t *testing.T) {
	cache := MkCache()
	producer := func() (Value, error) {
		return dao.GetById("123")
	}

	dao.Counter = 0
	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			for i := 0; i < 100; i++ {
				fv := cache.LoadOrElse("123", producer)
				value, _ := fv.Value()
				if data, ok := value.(MockedData); ok && data.id != "123" {
					t.Fatalf("Returned value is incorrect: %v", data)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if dao.Counter > 1 {
		t.Fatalf("dao has been invoked more than one time")
	}
}
