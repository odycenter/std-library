// Package crash 异常处理模块
package crash

import (
	"fmt"
	"reflect"
)

// CacheInf 异常处理模块抽象
type CacheInf interface {
	Catch(e error, fn func(err error)) CacheInf
	Finally
}

// Finally CacheInf 异常处理模块后续执行抽象
type Finally interface {
	Finally(fn ...func())
}

type cache struct {
	err  error
	done bool //是否已执行过cache
}

func (c *cache) cached() bool {
	if c.done || c.err == nil {
		return false
	}
	return true
}

// Catch 判断是否是指定类型的Error，如果e为nil，则默认处理全部err
func (c *cache) Catch(e error, fn func(err error)) CacheInf {
	if !c.cached() {
		return c
	}
	if e == nil || reflect.TypeOf(e) == reflect.TypeOf(c.err) {
		fn(c.err)
		c.done = true
	}
	return c
}

// Finally 遍历并在Finally函数执行完毕之后执行
func (c *cache) Finally(fn ...func()) {
	for _, f := range fn {
		defer f()
	}
	if c.err != nil && !c.done {
		panic(c.err)
	}
}

// Try 将要try cache的func传入执行
// Warning
// 程序稳定性有一定影响，会影响程序执行速度
// 请勿在for中调用
func Try(fn func()) CacheInf {
	c := new(cache)
	defer func() {
		defer func() {
			if e := recover(); e != nil {
				switch e.(type) {
				case string:
					c.err = fmt.Errorf(e.(string))
				default:
					c.err = e.(error)
				}
			}
		}()
		fn()
	}()
	return c
}
