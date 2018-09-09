package future

import "encoding/json"

type Future interface {
	Value() (interface{}, error)
	Then(func(interface{}, error) (interface{}, error)) Future
	OnSuccess(func(interface{})) Future
	OnFailure(func(error)) Future
}

func MakeFuture(producer func() (interface{}, error)) Future {
	fv := &futureImpl{
		done: make(chan struct{}),
	}
	go func() {
		fv.value, fv.err = producer()
		close(fv.done) //notify all the waiting/running Value() goroutine
	}()
	return fv
}

func MakeValue(value interface{}, err error) Future {
	done := make(chan struct{})
	defer close(done)
	return &futureImpl{
		value: value,
		err:   err,
		done:  done,
	}
}

type futureImpl struct {
	value interface{}
	err   error
	done  chan struct{} //ignore the value in chan
}

func (fv *futureImpl) Value() (interface{}, error) {
	<-fv.done //waiting for the result
	return fv.value, fv.err
}

func (fv *futureImpl) OnSuccess(f func(interface{})) Future {
	go func() {
		<-fv.done
		if fv.err == nil {
			f(fv.value)
		}
	}()
	return fv //returns a chained future
}

func (fv *futureImpl) OnFailure(f func(error)) Future {
	go func() {
		<-fv.done
		if fv.err != nil {
			f(fv.err)
		}
	}()
	return fv //returns a chained future
}

func (fv *futureImpl) Then(f func(interface{}, error) (interface{}, error)) Future {
	return MakeFuture(func() (interface{}, error) {
		<-fv.done
		v, err := f(fv.value, fv.err)
		return v, err
	})
}

func (fv *futureImpl) MarshalJSON() ([]byte, error) {
	value, _ := fv.Value()
	return json.Marshal(value)
}
