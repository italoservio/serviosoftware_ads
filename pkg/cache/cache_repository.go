package cache

import (
	"time"
)

type CacheRepository interface {
	Set(key string, value interface{})
	SetWithExpiration(key string, value interface{}, expiration time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
	Flush()
	ItemCount() int
}
