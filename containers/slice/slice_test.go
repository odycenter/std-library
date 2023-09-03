package slice_test

import (
	"fmt"
	"testing"

	"github.com/odycenter/std-library/containers/slice"
)

func TestMutex(t *testing.T) {
	m := slice.New[int](new(slice.Mutex[int]))
	m.Push(1)
	fmt.Println(m.Index(0))
	m.Push(2)
	m.Push(3)
	m.Delete(3)
	fmt.Println(m.Len())
	err := m.Range(func(v int) error {
		fmt.Println("V=>", v)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func TestRW(t *testing.T) {
	m := slice.New[int](new(slice.RWMutex[int]))
	m.Push(1)
	fmt.Println(m.Index(0))
	m.Push(2)
	m.Push(3)
	m.Delete(1)
	fmt.Println(m.Len())
	err := m.Range(func(v int) error {
		fmt.Println(v)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
