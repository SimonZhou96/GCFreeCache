package GeeCache

import (
	"GeeCache/lru"
	"fmt"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"
)

var (
	nameSpace *CacheNameSpace
)
func init() {
	nameSpace = NewCacheNameSpace(lru.NewCache(10e8, nil), "ben_mart")
}
func TestCacheNameSpace_Add(t *testing.T) {
	var wg sync.WaitGroup
	var go_routine_num = 100
	start := time.Now()
	for i := 0; i < go_routine_num; i++ {
		wg.Add(1)
		key := fmt.Sprintf("key%v", i)
		value := fmt.Sprintf("%v",i)
		unit := lru.NewCacheUnit([]byte(key), []byte(value))
		go func() {
			nameSpace.Add(unit)
			defer wg.Done()
		}()
	}
	wg.Wait()
	//
	log.Println(fmt.Sprintf("add time: %v", time.Since(start)))
	var wg1 sync.WaitGroup
	wg1.Add(go_routine_num)
	start = time.Now()
	for i := 0; i < go_routine_num; i++ {
		key := fmt.Sprintf("key%v", i)
		go func() {
			defer wg1.Done()
			value, ok := nameSpace.Get(key, nil)
			if ok {
				log.Println(value.(*lru.CacheUnit).GetStringValue())
			}
		}()
	}
	wg1.Wait()
	log.Println(fmt.Sprintf("get time: %v", time.Since(start)))
}

func TestCacheNameSpace_Delete(t *testing.T) {
	type fields struct {
		factory CacheMethodFactory
		label   string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &CacheNameSpace{
				factory: tt.fields.factory,
				label:   tt.fields.label,
			}
			if err := n.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCacheNameSpace_Get(t *testing.T) {
	type fields struct {
		factory CacheMethodFactory
		label   string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &CacheNameSpace{
				factory: tt.fields.factory,
				label:   tt.fields.label,
			}
			got, got1 := n.Get(tt.args.key, nil)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCacheNameSpace_Update(t *testing.T) {
	type fields struct {
		factory CacheMethodFactory
		label   string
	}
	type args struct {
		unit *lru.CacheUnit
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &CacheNameSpace{
				factory: tt.fields.factory,
				label:   tt.fields.label,
			}
			if err := n.Update(tt.args.unit); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}