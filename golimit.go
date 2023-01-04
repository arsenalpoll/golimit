package golimit

import (
	"sync"
)
// 1
// 2
type Limit struct {
	max       uint             //并发最大数量
	count     uint             //当前已有并发数
	isAddLock bool             //是否已锁定增加
	zeroChan  chan interface{} //为0时广播
	addLock   sync.Mutex
	dataLock  sync.Mutex
}

func GoLimit(max uint) *Limit {
	return &Limit{max: max, count: 0, isAddLock: false, zeroChan: nil}
}

//并发计数加1.若 计数>=max_num, 则阻塞,直到 计数<max_num
func (g *Limit) Add() {
	g.addLock.Lock()
	g.dataLock.Lock()

	g.count += 1

	if g.count < g.max { //未超并发时解锁,后续可以继续增加
		g.addLock.Unlock()
	} else { //已到最大并发数, 不解锁并标记. 等数量减少后解锁
		g.isAddLock = true
	}

	g.dataLock.Unlock()
}

// 并发计数减1
// 若计数<max_num, 可以使原阻塞的Add()快速解除阻塞

func (g *Limit) Done() {
	g.dataLock.Lock()

	g.count -= 1

	//解锁
	if g.isAddLock == true && g.count < g.max {
		g.isAddLock = false
		g.addLock.Unlock()
	}

	//0广播
	if g.count == 0 && g.zeroChan != nil {
		close(g.zeroChan)
		g.zeroChan = nil
	}

	g.dataLock.Unlock()
}

//更新最大并发计数为, 若是调大, 可以使原阻塞的Add()快速解除阻塞
func (g *Limit) SetMax(n uint) {
	g.dataLock.Lock()

	g.max = n

	//解锁
	if g.isAddLock == true && g.count < g.max {
		g.isAddLock = false
		g.addLock.Unlock()
	}

	//加锁
	if g.isAddLock == false && g.count >= g.max {
		g.isAddLock = true
		g.addLock.Lock()
	}

	g.dataLock.Unlock()
}

// 等待
func (g *Limit) Wait() {
	g.dataLock.Lock()

	if g.count == 0 {
		g.dataLock.Unlock()
		return
	}

	if g.zeroChan == nil {
		g.zeroChan = make(chan interface{})
	}

	c := g.zeroChan
	g.dataLock.Unlock()

	<-c
}

//获取并发计数
func (g *Limit) Count() uint {
	return g.count
}

//获取最大并发计数
func (g *Limit) Max() uint {
	return g.max
}
