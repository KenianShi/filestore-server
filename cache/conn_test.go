package cache

import (
	"testing"
	"fmt"
)

func TestRedisPool(t *testing.T) {
	pool := RedisPool()
	fmt.Println(pool)
}
