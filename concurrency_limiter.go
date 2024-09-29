package concurrencyLimiter

import (
	"sync"
)

type ConcurrencyLimiter struct {
	runningNum  int32
	limit       int32
	blockingNum int32
	cond        *sync.Cond
	mu          *sync.Mutex
}

// NewConcurrencyLimiter 创建一个并发限制器，limit为并发限制数量，可通过 Reset() 动态调整limit。
// 每次调用 Get() 来获取一个资源，然后创建一个协程，完成任务后通过 Release() 释放资源。
func NewConcurrencyLimiter(limit int32) *ConcurrencyLimiter {
	l := new(sync.Mutex)
	return &ConcurrencyLimiter{
		limit: limit,
		cond:  sync.NewCond(l),
		mu:    l,
	}
}

// Reset 可更新limit，需要保证limit > 0
func (c *ConcurrencyLimiter) Reset(limit int32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	tmp := c.limit
	c.limit = limit
	blockingNum := c.blockingNum
	// 优先唤醒阻塞的任务
	if limit-tmp > 0 && blockingNum > 0 {
		for i := int32(0); i < limit-tmp && blockingNum > 0; i++ {
			c.cond.Signal()
			blockingNum--
		}
	}
}

// Get 当 concurrencyLimiter 没有资源时，会阻塞。
func (c *ConcurrencyLimiter) Get() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.runningNum < c.limit {
		c.runningNum++
		return
	}
	c.blockingNum++
	for !(c.runningNum < c.limit) {
		c.cond.Wait()
	}
	c.runningNum++
	c.blockingNum--
}

// Release 释放一个资源
func (c *ConcurrencyLimiter) Release() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.blockingNum > 0 {
		c.runningNum--
		c.cond.Signal()
		return
	}

	c.runningNum--
}
