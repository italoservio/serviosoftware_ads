package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type GoCacheRepository struct {
	client *cache.Cache
}

func NewGoCacheRepository(defaultExpiration, cleanupInterval time.Duration) CacheRepository {
	return &GoCacheRepository{
		client: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *GoCacheRepository) Set(key string, value interface{}) {
	c.client.Set(key, value, cache.DefaultExpiration)
}

func (c *GoCacheRepository) SetWithExpiration(key string, value interface{}, expiration time.Duration) {
	c.client.Set(key, value, expiration)
}

func (c *GoCacheRepository) Get(key string) (interface{}, bool) {
	return c.client.Get(key)
}

func (c *GoCacheRepository) Delete(key string) {
	c.client.Delete(key)
}

func (c *GoCacheRepository) Flush() {
	c.client.Flush()
}

func (c *GoCacheRepository) ItemCount() int {
	return c.client.ItemCount()
}
