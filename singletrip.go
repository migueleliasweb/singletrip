package singletrip

import (
	"context"
	"sync/atomic"

	uberatomic "go.uber.org/atomic"
)

type SingleTrip[T any] struct {
	innerVal           atomic.Value
	innerErr           uberatomic.Error
	finishedPopulating chan struct{}
	shouldDo           atomic.Bool
}

func (st *SingleTrip[T]) Reset() {
	st.shouldDo.Store(true)
	st.innerErr.Store(nil)
	st.innerVal.Store(nil)
}

func (st *SingleTrip[T]) Do(
	ctx context.Context,
	key string,
	fn func(ctx context.Context) (T, error),
) (
	val T,
	err error,
) {
	// if we can swap it to false, then the last value was true
	if st.shouldDo.Swap(false) {
		val, err = fn(ctx)

		st.innerVal.Store(val)
		st.innerErr.Store(err)

		close(st.finishedPopulating)
	}

	if st.innerVal.Load() == nil {
		<-st.finishedPopulating
	}

	return st.innerVal.Load().(T), st.innerErr.Load()
}

func NewSingleTrip[T any]() *SingleTrip[T] {
	st := &SingleTrip[T]{
		innerVal:           atomic.Value{},
		innerErr:           uberatomic.Error{},
		shouldDo:           atomic.Bool{},
		finishedPopulating: make(chan struct{}),
	}

	st.shouldDo.Store(true)

	return st
}
