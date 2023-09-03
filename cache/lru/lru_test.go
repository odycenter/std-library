package lru_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/odycenter/std-library/cache/lru"
)

func TestLRU(t *testing.T) {
	cache := lru.New[int, string](10)
	for i := 0; i < 15; i++ {
		time.Sleep(time.Microsecond * 100)
		cache.Put(i, time.Now().Format(time.RFC3339Nano))
	}
	fmt.Println(cache.Get(0))
	fmt.Println(cache.Get(14))
	fmt.Println(cache.Get(5))
	cache.Resize(1)
	fmt.Println(cache.Get(10))
}

func TestLRUObject(t *testing.T) {
	type S struct {
		A string
	}
	cache := lru.New[int, S](10)
	for i := 0; i < 15; i++ {
		time.Sleep(time.Microsecond * 100)
		cache.Put(i, S{time.Now().Format(time.RFC3339Nano)})
	}
	fmt.Println(cache.Get(0))
	fmt.Println(cache.Get(14))
	fmt.Println(cache.Get(5))
	cache.Resize(1)
	fmt.Println(cache.Get(10))
}

func TestAsyncLRU(t *testing.T) {
	cache := lru.New[int, string](10)
	go func() {
		c := 0
		for t := range time.Tick(time.Millisecond) {
			c++
			cache.Put(c, t.Format(time.RFC3339Nano))
		}
	}()
	go func() {
		c := 0
		for range time.Tick(time.Millisecond * 2) {
			c++
			fmt.Println(cache.Get(c))
		}
	}()
	go func() {
		c := 0
		for range time.Tick(time.Millisecond * 5) {
			c++
			fmt.Println(cache.Get(c))
		}
	}()
	<-time.After(time.Hour)
}

func TestLRUObject1(t *testing.T) {
	var wg sync.WaitGroup
	type S struct {
		Account        string    `bson:"Account"`
		PlayerId       int32     `bson:"PlayerId"`
		AgentId        int32     `bson:"AgentId"`
		ChannelId      string    `bson:"ChannelId"`
		PackId         int32     `bson:"PackId"`
		ServerVersion  int32     `bson:"ServerVersion"`
		AuthToken      string    `bson:"AuthToken"`
		ExpireDate     int32     `bson:"ExpireDate"`
		Delete         bool      `bson:"Delete"`
		VendorPlayerId string    `bson:"VendorPlayerId"`
		FirstBetGold   int64     `bson:"FirstBetGold"`
		CreateTime     time.Time `bson:"CreateTime"`
		UpdateTime     time.Time `bson:"UpdateTime"`
	}
	cache := lru.New[string, S](10000)
	for i := 0; i < 20000; i++ {
		cache.Put(fmt.Sprint(i), S{Account: time.Now().Format(time.RFC3339Nano)})
	}
	for i := 10000; i < 20000; i++ {
		wg.Add(1)
		//time.Sleep(time.Microsecond * 100)
		go func(i int) {
			cache.Put(fmt.Sprint(i), S{Account: time.Now().Format(time.RFC3339Nano)})
			defer wg.Done()
		}(i)
	}
	for i := 1000; i < 50000; i++ {
		wg.Add(1)
		//time.Sleep(time.Microsecond * 100)
		go func(i int) {
			cache.Put(fmt.Sprint(i), S{Account: time.Now().Format(time.RFC3339Nano)})
			defer wg.Done()
		}(i)
	}
	for i := 1000; i < 50000; i++ {
		wg.Add(1)
		//time.Sleep(time.Microsecond * 100)
		go func(i int) {
			cache.Get(fmt.Sprint(i))
			defer wg.Done()
		}(i)
	}
	for i := 0; i < 20000; i++ {
		wg.Add(1)
		//time.Sleep(time.Microsecond * 100)
		go func(i int) {
			cache.Get(fmt.Sprint(i))
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println(cache.Get("0"))
	fmt.Println(cache.Get("14"))
	fmt.Println(cache.Get("5"))
	cache.Resize(1)
	fmt.Println(cache.Get("10"))
}
