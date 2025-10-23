package cache

import "time"

type Cache interface {
	Get(key string) (string, error) // return "", nil if MISS
	Set(key, val string, ttl time.Duration) error
	Del(keys ...string) error
}
