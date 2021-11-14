package g_set

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/proto"
)

func TestGSetAddContains(t *testing.T) {
	a, _ := New([]proto.Message{})

	Insert(a, &status.Status{Code: 1})
	Insert(a, &status.Status{Code: 1})

	containsOne, err := Contains(a, &status.Status{Code: 1})
	require.NoError(t, err)
	containsTwo, err := Contains(a, &status.Status{Code: 2})
	require.NoError(t, err)

	require.True(t, containsOne)
	require.False(t, containsTwo)

	elements, err := Elements(a)
	require.NoError(t, err)

	require.Len(t, elements, 1)
	require.True(t, proto.Equal(&status.Status{Code: 1}, elements[0]))

	Insert(a, &status.Status{Code: 2})

	// Insert 2
	containsTwo, err = Contains(a, &status.Status{Code: 2})
	require.NoError(t, err)
	require.True(t, containsTwo)
}

func TestGSetMerge(t *testing.T) {
	a, err := New([]proto.Message{
		&status.Status{Code: 1},
		&status.Status{Code: 2},
	})
	require.NoError(t, err)
	b, err := New([]proto.Message{
		&status.Status{Code: 2},
		&status.Status{Code: 3},
	})
	require.NoError(t, err)

	c, err := Merge(a, b)
	require.NoError(t, err)

	assert.Equal(t, 3, Len(c))

	for _, msg := range []proto.Message{
		&status.Status{Code: 1},
		&status.Status{Code: 2},
		&status.Status{Code: 3},
	} {
		contains, err := Contains(c, msg)
		require.NoError(t, err)

		assert.True(t, contains)
	}

	containsFour, err := Contains(c, &status.Status{Code: 4})
	require.NoError(t, err)
	assert.False(t, containsFour)
}
