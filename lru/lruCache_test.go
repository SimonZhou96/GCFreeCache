package lru

import (
	"log"
	"reflect"
	"testing"
)

var (
	cache *LRUCache
)
func init()  {
	cache = NewCache(100, nil)
}
func TestNewCache(t *testing.T) {
	type args struct {
		maxBytes  int64
		onEvicted func(string, *CacheUnit)
	}
	tests := []struct {
		name string
		args args
		want *LRUCache
	}{
		{
			name: "",
			args: args{
				maxBytes:  100,
				onEvicted: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCache(tt.args.maxBytes, tt.args.onEvicted); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Get(t *testing.T) {
	cache.Get("hello")
}

func TestCache_Add(t *testing.T) {
	cache.Add(NewCacheUnit(
		[]byte("hello"),
		[]byte("world"),
	))
}

func TestCache_Update(t *testing.T) {
	cache.Add(NewCacheUnit(
		[]byte("hello"),
		[]byte("world"),
	))
	value, _ := cache.Get("hello")
	log.Println(value)
	cache.Update(NewCacheUnit(
		[]byte("hello"),
		[]byte("world22w"),
	))
	value, _ = cache.Get("hello")
	log.Println(value)
	cache.Delete("hello")
	value, _ = cache.Get("hello")
	log.Println(value)
}