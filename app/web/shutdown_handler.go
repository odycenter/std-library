package web

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type ShutdownHandler struct {
	activeRequests int32
	cond           *sync.Cond
	shutdown       int32 // 0 for false, 1 for true
}

func NewShutdownHandler() *ShutdownHandler {
	return &ShutdownHandler{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (s *ShutdownHandler) Increment() {
	atomic.AddInt32(&s.activeRequests, 1)
}

func (s *ShutdownHandler) Decrement() {
	atomic.AddInt32(&s.activeRequests, -1)
	s.cond.Signal()
}

func (s *ShutdownHandler) IsShutdown() bool {
	return atomic.LoadInt32(&s.shutdown) == 1
}

func (s *ShutdownHandler) Shutdown() {
	atomic.StoreInt32(&s.shutdown, 1) // set shutdown to true
}

func (s *ShutdownHandler) ActiveRequests() int32 {
	return atomic.LoadInt32(&s.activeRequests)
}

func (s *ShutdownHandler) AwaitTermination(timeoutInMs int64) bool {
	innerCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInMs)*time.Millisecond)
	defer cancel()

	done := make(chan bool, 1)
	go func() {
		done <- s.awaitTermination(timeoutInMs)
	}()

	select {
	case <-innerCtx.Done():
		return false
	case success := <-done:
		return success
	}
}

func (s *ShutdownHandler) awaitTermination(timeoutInMs int64) bool {
	deadline := time.Now().Add(time.Duration(timeoutInMs) * time.Millisecond)
	s.cond.L.Lock()
	for s.ActiveRequests() > 0 {
		timeLeft := time.Until(deadline)
		if timeLeft <= 0 {
			s.cond.L.Unlock()
			return false
		}
		s.cond.Wait()
	}
	s.cond.L.Unlock()
	return true
}
