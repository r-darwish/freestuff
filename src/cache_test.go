package freestuff

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRedis(t *testing.T) {
	cache := NewRedisCache()

	key := "something"

	_, err := cache.conn.Do("DEL", key)
	require.NoError(t, err)

	known, err := cache.IsKnown(key)
	assert.False(t, known)
	assert.NoError(t, err)

	err = cache.SetKnown(key)
	assert.NoError(t, err)

	known, err = cache.IsKnown(key)
	assert.True(t, known)
	assert.NoError(t, err)
}
