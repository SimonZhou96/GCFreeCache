package lru

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"log"
	"sync"
)

const (
	OffsetKeyByteSize = 4
	OffSetValueByteSize = 4
)

type LRUCache struct {
	rwLock sync.RWMutex
	maxBytes int64
	nbytes int64
	ll *list.List
	cache sync.Map
	OnEvicted func(key string, value *CacheUnit)
}

type CacheUnit struct {
	//element []byte  // value offset长度+value的形式
	Key []byte
	Value []byte
}

func NewCache(maxBytes int64, onEvicted func(string, *CacheUnit)) *LRUCache {
	return &LRUCache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		OnEvicted: onEvicted,
	}
}

func Int32ToByteArray(i int32) []byte {
	buf := new(bytes.Buffer)
	var num = i
	// 采用bigendian 策略
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	keyOffsetBytes := buf.Bytes()
	return keyOffsetBytes
}

func NewCacheUnit(key []byte, value []byte) *CacheUnit {

	keyOffsetBytes := Int32ToByteArray(int32(len(key)))
	valueOffsetBytes:= Int32ToByteArray(int32(len(value)))
	return &CacheUnit{
		Key: append(keyOffsetBytes, key...),
		Value: append(valueOffsetBytes, value...),
	}
}

func (cu *CacheUnit) GetKey() []byte {
	if cu == nil {
		return nil
	}
	offsetBytes := cu.Key[:OffsetKeyByteSize]
	offset := binary.BigEndian.Uint32(offsetBytes)
	key := cu.Key[OffsetKeyByteSize : OffsetKeyByteSize+ offset]
	return key
}

func (cu *CacheUnit) GetValue() []byte {
	if cu == nil {
		return nil
	}
	valueoffset := binary.BigEndian.Uint32(cu.Value[:OffSetValueByteSize])
	value := cu.Value[OffSetValueByteSize :valueoffset+OffSetValueByteSize]
	return value
}

func (cu *CacheUnit) TotalLen() int64 {
	return cu.ValueLen() + cu.KeyLen()
}

func (cu *CacheUnit) ValueLen() int64 {
	return int64(len(cu.GetValue())) + OffSetValueByteSize
}

func (cu *CacheUnit) KeyLen() int64 {
	return int64(len(cu.GetKey())) + OffsetKeyByteSize
}

func (cu *CacheUnit) GetStringKey() string {
	return string(cu.GetKey())
}

func (cu *CacheUnit) GetStringValue() string {
	return string(cu.GetValue())
}

func (cu *LRUCache) EntryLen() int64 {
	return int64(cu.ll.Len())
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	value, ok := c.cache.Load(key)
	if !ok {
		return nil, false
	}
	// lru operation
	c.ll.MoveToFront(value.(*list.Element))
	return value.(*list.Element).Value.(*CacheUnit), ok
}

// RemoveOldest removes the oldest item
func (c *LRUCache) removeOldest() {
	ele := c.ll.Back()
	if ele == nil {
		return
	}
	c.ll.Remove(ele)
	kv := ele.Value.(*CacheUnit)
	c.cache.Delete(kv.GetStringKey())
	c.nbytes -= kv.TotalLen()
	if c.OnEvicted != nil {
		c.OnEvicted(kv.GetStringKey(), kv)
	}
}

func (c *LRUCache) Add(unit *CacheUnit) {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	if unit == nil {
		return
	}
	if unit.TotalLen() > c.maxBytes {
		log.Fatal(fmt.Sprintf("overflow! new unit size: %v, capacity: %v", unit.TotalLen(), c.maxBytes))
		return
	}
	if ele, ok := c.cache.Load(unit.GetStringKey()); ok {
		ele := ele.(*list.Element)
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*CacheUnit)
		c.nbytes += unit.ValueLen() - kv.ValueLen()
		kv = unit
	} else {
		e := c.ll.PushFront(unit)
		c.cache.Store(unit.GetStringKey(), e)
		c.nbytes += unit.TotalLen()
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.removeOldest()
	}
}


func (c *LRUCache) Update(unit *CacheUnit) {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	if unit == nil {
		return
	}
	if unit.TotalLen() > c.maxBytes {
		log.Fatal(fmt.Sprintf("overflow! new unit size: %v, capacity: %v", unit.TotalLen(), c.maxBytes))
		return
	}
	if _, ok := c.cache.Load(unit.GetStringKey()); !ok {
		c.Add(unit)
		return
	}
	ele, _ := c.cache.Load(unit.GetStringKey())
	c.nbytes += unit.TotalLen() - ele.(*list.Element).Value.(*CacheUnit).TotalLen()
	if c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.removeOldest()
	}
	ele.(*list.Element).Value = unit
}

func (c *LRUCache) Delete(key string) {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	ele, ok := c.cache.Load(key)
	if !ok {
		return
	}
	c.ll.Remove(ele.(*list.Element))
	c.cache.Delete(key)
}