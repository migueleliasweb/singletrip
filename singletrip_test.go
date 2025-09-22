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
	total := 10
	callbackCalledNTimes := &atomic.Int32{}
	wg := sync.WaitGroup{}
	wg.Add(total)

	st := singletrip.NewSingleTrip[bool]()

	for range total {
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

func TestSingleTripAndReset(t *testing.T) {
	callbackCalledNTimes := &atomic.Int32{}
	wg := sync.WaitGroup{}
	total := 1000
	wg.Add(total)

	st := &singletrip.SingleTrip[bool]{}

	for i := range total {
		go func() {
			defer func() {
				wg.Done()
			}()

			if (total / 2) == i {
				st.Reset()
			}

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
		int32(2),
		callbackCalledNTimes.Load(),
	)
}
