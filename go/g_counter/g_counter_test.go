package g_counter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGCounter(t *testing.T) {
	tests := []struct {
		aIncs  []uint64
		bIncs  []uint64
		result int64
	}{
		{
			aIncs:  []uint64{10, 50, 100},
			bIncs:  []uint64{15, 55, 200},
			result: 430,
		},
		{
			aIncs:  []uint64{1, 2, 3},
			bIncs:  []uint64{},
			result: 6,
		},
		{
			aIncs:  []uint64{},
			bIncs:  []uint64{1, 2, 3},
			result: 6,
		},
		{
			aIncs:  []uint64{50},
			bIncs:  []uint64{100},
			result: 150,
		},
	}

	for _, test := range tests {
		a := New("a")

		for _, inc := range test.aIncs {
			Increment(a, inc)
		}

		b := New("b")
		for _, inc := range test.bIncs {
			Increment(b, inc)
		}

		require.Equal(t, test.result, Count(Merge("a", a, b)))
	}
}
