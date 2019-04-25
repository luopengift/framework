package limit

import (
	"fmt"
	"os"
	"sync"
)

// Limit channel
type Limit struct {
	wg      *sync.WaitGroup
	ch      chan []byte
	stop    chan struct{} // stop chan
	LogFunc func(string, ...interface{})
}

// NewLimit channel
func NewLimit(max int) *Limit {
	return &Limit{
		wg:   new(sync.WaitGroup),
		ch:   make(chan []byte, max),
		stop: make(chan struct{}),
		LogFunc: func(format string, v ...interface{}) {
			fmt.Fprintf(os.Stderr, format+"\n", v)
		},
	}
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

// Run run
func (c *Limit) Run(fun func() error) error {
	c.Add()
	go func() {
		defer c.Done()
		if err := fun(); err != nil {
			c.LogFunc("%s", err)
		}
	}()
	return nil
}
