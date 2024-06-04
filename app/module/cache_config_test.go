package module_test

import (
	"context"
	"fmt"
	"std-library/app/cache"
	"std-library/app/module"
	"testing"
	"time"
)

type CacheTest struct {
	PlayerName string
}
type CacheTest2 struct {
	PlayerName string
}

func TestCacheNoConnect(t *testing.T) {
	moduleContext := &module.Context{}
	moduleContext.Initialize()
	config := module.CacheConfig{}
	config.Initialize(moduleContext, "test")
	config.Redis("redis")
	c := config.Add(CacheTest{}, time.Second*300)

	key := "test"
	c.Evict(context.Background(), key)
	c.Evict(context.Background(), "test2")
	values := map[string]interface{}{key: CacheTest{PlayerName: "data1 from put all"}, "test2": CacheTest{PlayerName: "data2 from put all"}}
	c.PutAll(context.Background(), values)
	get(c, key)
	get(c, "test2")

	result, err := c.GetAll(context.Background(), []string{key, "test2"}, CacheTest{}, func(key string) (interface{}, error) {
		return CacheTest{PlayerName: "from get all loader:" + key}, nil
		//return nil, nil
	})
	if err != nil {
		panic(err)
	}
	for k, v := range result {
		if t, ok := (v).(CacheTest); ok {
			fmt.Printf("key: %s, value:%v \r\n", k, t.PlayerName)
		}
	}

	//c.Put(context.Background(), key, CacheTest{PlayerName: "from put"})
	//get(c, key)
	//
	//for i := 0; i < 100; i++ {
	//	get(c, key)
	//	// sleep 1 second
	//	time.Sleep(time.Millisecond * 1000)
	//}

	// assert.Equal(t, "test:create from loader", obj.PlayerName)
}

func TestCacheMiss(t *testing.T) {
	moduleContext := &module.Context{}
	moduleContext.Initialize()
	config := module.CacheConfig{}
	config.Initialize(moduleContext, "test")
	config.Redis("mggroup-dev.mijptn.ng.0001.apne1.cache.amazonaws.com:6379")
	c := config.Add(CacheTest{}, time.Second*300)

	key := "test"
	c.Evict(context.Background(), key)
	c.Evict(context.Background(), "test2")
	values := map[string]interface{}{key: CacheTest{PlayerName: "data1 from put all"}, "test2": CacheTest{PlayerName: "data2 from put all"}}
	c.PutAll(context.Background(), values)
	get(c, key)
	get(c, "test2")

	result, err := c.GetAll(context.Background(), []string{key, "test2"}, CacheTest{}, func(key string) (interface{}, error) {
		return CacheTest{PlayerName: "from get all loader:" + key}, nil
		//return nil, nil
	})
	if err != nil {
		panic(err)
	}
	for k, v := range result {
		if t, ok := (v).(CacheTest); ok {
			fmt.Printf("key: %s, value:%v \r\n", k, t.PlayerName)
		}
	}

	//c.Put(context.Background(), key, CacheTest{PlayerName: "from put"})
	//get(c, key)
	//
	//for i := 0; i < 100; i++ {
	//	get(c, key)
	//	// sleep 1 second
	//	time.Sleep(time.Millisecond * 1000)
	//}

	// assert.Equal(t, "test:create from loader", obj.PlayerName)
}

func get(cache cache.Cache, key string) {
	var obj CacheTest
	err := cache.Get(context.Background(), key, &obj, func(key string) (interface{}, error) {
		println("call loader")
		return &CacheTest{PlayerName: key + ":create from loader"}, nil
	})
	if err != nil {
		panic(err)
	}
	println(">>>>>>>>" + obj.PlayerName)
}
