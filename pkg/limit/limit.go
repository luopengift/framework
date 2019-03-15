package limit

import (
	"fmt"
	"sync"

	"github.com/luopengift/log"
)

// Limit channel
type Limit struct {
	wg   *sync.WaitGroup
	ch   chan []byte
	stop chan struct{} // stop chan
}

// NewLimit channel
func NewLimit(max int) *Limit {
	return &Limit{
		wg:   new(sync.WaitGroup),
		ch:   make(chan []byte, max),
		stop: make(chan struct{}),
	}
}

// Close channel
func (c *Limit) Close() error {
	// for i := 0; i < 10; i++ {
	// 	if len(c.ch) == 0 {
	// 		close(c)
	// 		return nil
	// 	}
	// 	time.Sleep(10 * time.Millisecond)
	// }
	return fmt.Errorf("closed ch failed! ch is not empty, len is %d", c.Len())
}

//Put 往管道中写数据
func (c *Limit) Put(b []byte) error {
	c.wg.Add(1)
	c.ch <- b
	return nil
}

// Get 从管道中读数据
func (c *Limit) Get() ([]byte, bool) {
	v, ok := <-c.ch
	c.wg.Done()
	return v, ok
}

// Wait 等待所有gor退出
func (c *Limit) Wait() {
	c.wg.Wait()
}

// Add 往管道中放入一个标记，记录活跃数值
func (c *Limit) Add() {
	c.Put([]byte{})
}

// Done 从管道中取出一个标记，减少活跃数值
func (c *Limit) Done() {
	c.Get()
}

// Cap cap
func (c *Limit) Cap() int { return cap(c.ch) }

// Len len
func (c *Limit) Len() int { return len(c.ch) }

// Run run
func (c *Limit) Run(fun func() error) error {
	c.Add()
	go func() {
		if err := fun(); err != nil {
			log.Error("%s", err)
		}
		c.Done()
	}()
	return nil
}
