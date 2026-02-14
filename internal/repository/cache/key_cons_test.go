package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessTokenKey(t *testing.T) {
	assert.Equal(t, "access:abc123", AccessTokenKey("abc123"))
}

func TestRefreshTokenKey(t *testing.T) {
	assert.Equal(t, "refresh:xyz789", RefreshTokenKey("xyz789"))
}

func TestTaskCacheKey(t *testing.T) {
	assert.Equal(t, "task:42", TaskCacheKey(42))
}

func TestTaskListCacheKeyWithParams(t *testing.T) {
	result := TaskListCacheKeyWithParams(20, 0, "Created")
	assert.Equal(t, "tasks:list:20:0:Created", result)
}
