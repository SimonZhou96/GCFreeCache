package GeeCache

import (
	"GeeCache/lru"
)

type CacheMethodFactory interface {
	Get(key string) (interface{}, bool)
	Add(unit *lru.CacheUnit)
	Update(unit *lru.CacheUnit)
	Delete(key string)
}


