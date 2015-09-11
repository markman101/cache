package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	sync.RWMutex
	Key       interface{}
	Data      interface{}
	LifeSpan  time.Duration
	TimeStamp time.Time
}

func CreateCacheItem(cache_key interface{}, cache_data interface{}) *CacheItem {
	Item := CacheItem{Key: cache_key, Data: cache_data}
	return &Item
}
