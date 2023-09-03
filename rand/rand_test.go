package rand_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/odycenter/std-library/rand"
)

func TestRand(t *testing.T) {
	r := rand.Rand()
	for i := 0; i < 10; i++ {
		fmt.Println(r.Range(1, 10))
	}
	fmt.Println(rand.Rand().Range(1, 10))
	fmt.Println(rand.Rand().Range(1, 10))
	fmt.Println(rand.Rand().Range(1, 10))
	fmt.Println(rand.Rand().Range(1, 10))
	fmt.Println(rand.Rand().Range(1, 10))
	fmt.Println(rand.Rand().Range(1, 10))
	fmt.Println(rand.Rand().Strings(10))
	fmt.Println(rand.Rand().Number(10))
	fmt.Println(rand.Rand().General(10))
	fmt.Println(rand.Rand().Custom(10, "testDictionary"))
	fmt.Println(rand.Rand().Letters(10))
}

func TestRandWithSeed(t *testing.T) {
	r := rand.Rand(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		fmt.Println(r.Range(1, 10))
	}
	fmt.Println(r.Range(1, 10))
	fmt.Println(r.Range(1, 10))
	fmt.Println(r.Range(1, 10))
	fmt.Println(r.Range(1, 10))
	fmt.Println(r.Range(1, 10))
	fmt.Println(r.Range(1, 10))
	fmt.Println(r.Strings(10))
	fmt.Println(r.General(10))
	fmt.Println(r.Number(10))
	fmt.Println(r.Custom(10, "testDictionary"))
	fmt.Println(r.Letters(10))
}

func TestRandAsync(t *testing.T) {
	go randAsync()
	go randAsync()
	go randAsync()
	go randAsync()
	<-time.After(time.Second * 10)
}

func randAsync() {
	for range time.Tick(time.Millisecond) {
		fmt.Println(rand.Rand().Range(1, 10))
	}
}

func TestRandAsyncWithSeed(t *testing.T) {
	r := rand.Rand(time.Now().UnixNano())
	go randAsyncWithSeed(r)
	go randAsyncWithSeed(r)
	go randAsyncWithSeed(r)
	go randAsyncWithSeed(r)
	<-time.After(time.Second * 10)
}

func randAsyncWithSeed(r *rand.R) {
	for range time.Tick(time.Millisecond) {
		fmt.Println(r.Range(1, 10))
	}
}
