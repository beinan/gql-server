package future

import (
	"context"
	"encoding/json"
)

type Value = interface{}

type Future interface {
	Value() (Value, error)
	CancellableValue(context.Context) (Value, error)
	Then(func(Value) (Value, error)) Future
	OnSuccess(func(Value)) Future
	OnFailure(func(error)) Future
}

func MakeFuture(producer func() (Value, error)) Future {
	fv := &futureImpl{
		done: make(chan struct{}),
	}
	go func() {
		fv.value, fv.err = producer()
		close(fv.done) //notify all the waiting/running Value()
	}()
	return fv
}

func MakeValue(value Value, err error) Future {
	done := make(chan struct{})
	defer close(done)
	return &futureImpl{
		value: value,
		err:   err,
		done:  done,
	}
}

type futureImpl struct {
	value Value
	err   error
	done  chan struct{} //ignore the value in chan
}

func (fv *futureImpl) Value() (Value, error) {
	<-fv.done //waiting for the result
	return fv.value, fv.err
}

func (fv *futureImpl) CancellableValue(ctx context.Context) (Value, error) {
	select {
	case <-fv.done:
		return fv.value, fv.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (fv *futureImpl) OnSuccess(f func(Value)) Future {
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

func (fv *futureImpl) Then(f func(Value) (Value, error)) Future {
	return MakeFuture(func() (Value, error) {
		<-fv.done
		if fv.err != nil {
			return nil, fv.err
		}
		v, err := f(fv.value)
		return v, err
	})
}

func (fv *futureImpl) MarshalJSON() ([]byte, error) {
	value, _ := fv.Value()
	return json.Marshal(value)
}
