package singletrip_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/migueleliasweb/singletrip"
	"github.com/stretchr/testify/require"
)

func TestSingleTrip(t *testing.T) {
	callbackCalledNTimes := &atomic.Int32{}
	wg := sync.WaitGroup{}
	wg.Add(1000)

	st := &singletrip.SingleTrip[bool]{}

	for range 1000 {
		go func() {
			defer func() {
				wg.Done()
			}()

			res, err := st.Do(
				t.Context(),
				"something",
				func(ctx context.Context) (bool, error) {
					callbackCalledNTimes.Add(1)
					time.Sleep(time.Second * 5)
					return true, nil
				},
			)

			require.True(t, res)
			require.NoError(t, err)
		}()
	}

	wg.Wait()

	require.Equal(
		t,
		int32(1),
		callbackCalledNTimes.Load(),
	)
}
