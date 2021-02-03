package GeeCache

import (
	"GeeCache/lru"
	"fmt"
)

type getMissCallBack func(key string) (interface{}, bool)
// NameSpace 用户每一种请求都有一个命名空间，逻辑上的划分
type CacheNameSpace struct {
	factory CacheMethodFactory
	label   string
}

// NewCacheNameSpace 初始化新的命名空间
func NewCacheNameSpace(factory CacheMethodFactory, label string) *CacheNameSpace {
	return &CacheNameSpace{
		factory: factory,
		label:   label,
	}
}

func (n *CacheNameSpace) Add(unit *lru.CacheUnit) error {
	if n == nil {
		return fmt.Errorf("name space doesn't exist!")
	}
	n.factory.Add(unit)
	return nil
}

func (n *CacheNameSpace) Update(unit *lru.CacheUnit) error {
	if n == nil {
		return fmt.Errorf("name space doesn't exist!")
	}
	n.factory.Update(unit)
	return nil
}

func (n *CacheNameSpace) Delete(key string) error {
	if n == nil {
		return fmt.Errorf("name space doesn't exist!")
	}
	n.factory.Delete(key)
	return nil
}

// 回源函数
func (n *CacheNameSpace) Get(key string, callback getMissCallBack) (interface{}, bool) {
	if n == nil {
		return nil, false
	}
	value, ok := n.factory.Get(key)
	if !ok {
		// 回源
		return callback(key)
	}
	return value, true
}

