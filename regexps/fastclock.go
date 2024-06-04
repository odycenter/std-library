package regexps

import (
	"sync"
	"sync/atomic"
	"time"
)

// fastTime holds a time value (ticks since clock initialization)
type fastTime int64

// fastClock 提供快速时钟实现。
//
// 后台 goroutine 定期存储当前时间
// 到一个原子变量中。
//
// 通过比较可以快速检查截止日期是否过期
// 它的值存储在原子变量中的时钟。
//
// 一旦到达clockEnd，goroutine就会自动停止。
// (clockEnd 涵盖了迄今为止看到的最大截止日期 + 一些
// 额外的时间）。这确保了如果正则表达式超时
// 停止使用我们将停止后台工作。
type fastClock struct {
	// instances of atomicTime must be at the start of the struct (or at least 64-bit aligned)
	// otherwise 32-bit architectures will panic

	current  atomicTime // Current time (approximate)
	clockEnd atomicTime // When clock updater is supposed to stop (>= any existing deadline)

	// current and clockEnd can be read via atomic loads.
	// Reads and writes of other fields require mu to be held.
	mu      sync.Mutex
	start   time.Time // Time corresponding to fastTime(0)
	running bool      // Is a clock updater running?
}

var fast fastClock

// reached returns true if current time is at or past t.
func (t fastTime) reached() bool {
	return fast.current.read() >= t
}

// makeDeadline returns a time that is approximately time.Now().Add(d)
func makeDeadline(d time.Duration) fastTime {
	// Increase the deadline since the clock we are reading may be
	// just about to tick forwards.
	end := fast.current.read() + durationToTicks(d+clockPeriod)

	// Start or extend clock if necessary.
	if end > fast.clockEnd.read() {
		extendClock(end)
	}
	return end
}

// extendClock ensures that clock is live and will run until at least end.
func extendClock(end fastTime) {
	fast.mu.Lock()
	defer fast.mu.Unlock()

	if fast.start.IsZero() {
		fast.start = time.Now()
	}

	// Extend the running time to cover end as well as a bit of slop.
	if shutdown := end + durationToTicks(time.Second); shutdown > fast.clockEnd.read() {
		fast.clockEnd.write(shutdown)
	}

	// Start clock if necessary
	if !fast.running {
		fast.running = true
		go runClock()
	}
}

// stop the timeout clock in the background
// should only used for unit tests to abandon the background goroutine
func stopClock() {
	fast.mu.Lock()
	if fast.running {
		fast.clockEnd.write(fastTime(0))
	}
	fast.mu.Unlock()

	// pause until not running
	// get and release the lock
	isRunning := true
	for isRunning {
		time.Sleep(clockPeriod / 2)
		fast.mu.Lock()
		isRunning = fast.running
		fast.mu.Unlock()
	}
}

func durationToTicks(d time.Duration) fastTime {
	// Downscale nanoseconds to approximately a millisecond so that we can avoid
	// overflow even if the caller passes in math.MaxInt64.
	return fastTime(d) >> 20
}

const DefaultClockPeriod = 100 * time.Millisecond

// clockPeriod is the approximate interval between updates of approximateClock.
var clockPeriod = DefaultClockPeriod

func runClock() {
	fast.mu.Lock()
	defer fast.mu.Unlock()

	for fast.current.read() <= fast.clockEnd.read() {
		// Unlock while sleeping.
		fast.mu.Unlock()
		time.Sleep(clockPeriod)
		fast.mu.Lock()

		newTime := durationToTicks(time.Since(fast.start))
		fast.current.write(newTime)
	}
	fast.running = false
}

type atomicTime struct{ v int64 } // Should change to atomic.Int64 when we can use go 1.19

func (t *atomicTime) read() fastTime   { return fastTime(atomic.LoadInt64(&t.v)) }
func (t *atomicTime) write(v fastTime) { atomic.StoreInt64(&t.v, int64(v)) }
