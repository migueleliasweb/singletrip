package singletrip

import (
	"context"
	"sync"
)

type SingleTrip[T any] struct {
	innerVal T
	innerErr error
	finished chan struct{}
	once     sync.Once
}

func (st *SingleTrip[T]) Do(
	ctx context.Context,
	key string,
	fn func(ctx context.Context) (T, error),
) (
	val T,
	err error,
) {
	st.once.Do(func() {
		st.innerVal, st.innerErr = fn(ctx)

		st.finished = make(chan struct{})

		close(st.finished)
	})

	<-st.finished

	return st.innerVal, st.innerErr
}
