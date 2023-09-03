package lfu_test

import (
	"fmt"
	"std-library/cache/lfu"
	"testing"
	"time"
)

func TestLFU(t *testing.T) {
	cache := lfu.New[int, string](10)
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

func TestAsyncLFU(t *testing.T) {
	cache := lfu.New[int, string](10)
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
