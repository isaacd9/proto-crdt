package pn_counter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPNCounter(t *testing.T) {
	tests := []struct {
		aIncs, aDecs []uint64
		bIncs, bDecs []uint64
		result       int64
	}{
		{aIncs: []uint64{200}, aDecs: []uint64{50},
			bIncs: []uint64{5}, bDecs: []uint64{3},
			result: 152},
	}

	for _, test := range tests {
		a := New("a")

		for _, inc := range test.aIncs {
			Increment(a, inc)
		}
		for _, dec := range test.aDecs {
			Decrement(a, dec)
		}

		b := New("b")
		for _, inc := range test.bIncs {
			Increment(b, inc)
		}
		for _, dec := range test.bDecs {
			Decrement(a, dec)
		}

		require.Equal(t, test.result, Value(Merge("a", a, b)))
	}
}
