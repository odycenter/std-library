package slice_test

import (
	"fmt"
	"github.com/odycenter/std-library/containers/slice"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestArrayPartition(t *testing.T) {
	var source []string
	for i := 'a'; i <= 'z'; i++ {
		source = append(source, string(i))
	}
	arr := slice.ArrayPartition(source, 7)
	assert.Equal(t, 4, len(arr))
	assert.EqualValues(t, []string{"a", "b", "c", "d", "e", "f", "g"}, arr[0])
	assert.EqualValues(t, []string{"v", "w", "x", "y", "z"}, arr[3])
	assert.EqualValues(t, [][]string{source}, slice.ArrayPartition(source, -1))

	assert.Equal(t, 9, len(slice.ArrayPartition(source, 3)))
	assert.Equal(t, 1, len(slice.ArrayPartition(source, 26)))
	assert.Equal(t, 1, len(slice.ArrayPartition(source, 30)))

}

func TestToMap(t *testing.T) {
	// test1
	ids := []int{1, 2}
	m := slice.ToMap(ids, func(vl int) string {
		return fmt.Sprintf("%d_key", vl)
	})
	assert.Equal(t, m["1_key"], 1)
	assert.Equal(t, m["2_key"], 2)
	// test2
	names := []string{"Taipei", "NewYork", "London"}
	mn := slice.ToMap(names, func(vl string) string {
		return vl + "_address"
	})
	assert.Equal(t, mn["Taipei_address"], "Taipei")
	assert.Equal(t, mn["NewYork_address"], "NewYork")
	assert.Equal(t, mn["London_address"], "London")
}
