package singletrip

import (
	"context"
	"sync"
	"sync/atomic"
)

type SingleTrip[T any] struct {
	innerVal              T
	innerErr              error
	finishedBootstrapping chan struct{}
	shouldDo              atomic.Bool
	mut                   sync.Mutex
	once                  sync.Once
}

func (st *SingleTrip[T]) Reset() {
	st.shouldDo.Store(true)
}

func (st *SingleTrip[T]) Do(
	ctx context.Context,
	key string,
	fn func(ctx context.Context) (T, error),
) (
	val T,
	err error,
) {
	defer st.mut.Unlock()

	// if we can swap it to false, then the last value was true
	if st.shouldDo.Swap(false) {
		st.mut.Lock()
		st.innerVal, st.innerErr = fn(ctx)

		st.once.Do(func() {
			close(st.finishedBootstrapping)
		})

		return st.innerVal, st.innerErr
	}

	<-st.finishedBootstrapping

	st.mut.Lock()
	return st.innerVal, st.innerErr
}

func NewSingleTrip[T any]() *SingleTrip[T] {
	st := &SingleTrip[T]{
		finishedBootstrapping: make(chan struct{}),
	}

	st.shouldDo.Store(true)

	return st
}
