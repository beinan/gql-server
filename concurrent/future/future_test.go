package future

import (
	"sync"
	"testing"
	"time"
)

func generateProducer(t *testing.T, value interface{}, err error) func() (interface{}, error) {
	return func() (interface{}, error) {
		t.Log("Producing value")
		defer t.Log("Value produced")
		time.Sleep(time.Duration(500 * time.Millisecond))
		return value, err
	}
}

func TestFutureValue(t *testing.T) {
	producer := generateProducer(t, 3, nil)
	fv := MakeFuture(producer)
	result, err := fv.Value()
	if result != 3 {
		t.Fatalf("Returned value is incorrect: %v", result)
	}
	if err != nil {
		t.Fatalf("Returned error is not nil: %v", err)
	}
}

func TestFutureValueConcurrent(t *testing.T) {
	producer := generateProducer(t, 5, nil)
	fv := MakeFuture(producer)

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			for i := 0; i < 100; i++ {
				value, err := fv.Value()
				if value != 5 {
					t.Fatalf("Returned value is incorrect: %v", value)
				}
				if err != nil {
					t.Fatalf("Returned error is not nil: %v", err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	//test future value again after all the coroutine finished
	value, err := fv.Value()
	if value != 5 {
		t.Fatalf("Returned value is incorrect: %v", value)
	}
	if err != nil {
		t.Fatalf("Returned error is not nil: %v", err)
	}

}
