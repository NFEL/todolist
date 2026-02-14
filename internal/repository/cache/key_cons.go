package cache

import "fmt"

const (
	AccessTokenPrefix  = "access:"
	RefreshTokenPrefix = "refresh:"
	TaskCachePrefix    = "task:"
	TaskListCacheKey   = "tasks:list"
)

func AccessTokenKey(jti string) string {
	return AccessTokenPrefix + jti
}

func RefreshTokenKey(jti string) string {
	return RefreshTokenPrefix + jti
}

func TaskCacheKey(id uint) string {
	return fmt.Sprintf("%s%d", TaskCachePrefix, id)
}

func TaskListCacheKeyWithParams(limit, offset int, status string) string {
	return fmt.Sprintf("%s:%d:%d:%s", TaskListCacheKey, limit, offset, status)
}
