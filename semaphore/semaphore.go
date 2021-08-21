package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Semaphore struct{
	max int64
	cur int64
	ch  chan struct{}
	w   *waiters
	m  sync.Mutex
}

type waiters struct{
	value int64
	c   chan struct{}
	next interface{}
}

func NewSemaphore(cap int) *Semaphore {
	ch:=make(chan struct{},cap)
	for i:=0;i<cap;i++{
		ch<- struct{}{}
	}
	return &Semaphore{max: int64(cap),cur:int64(cap),ch:ch}
}

func (s *Semaphore) take(n int){
	for i := 0; i < n; i++ {
		<-s.ch
		fmt.Println("读取值... n==",n)
	}
	return
}

func (s *Semaphore) send(n int){
	for i := 0; i < n; i++ {
		fmt.Println("插入... n==",n)
		s.ch<- struct{}{}
	}
	return
}

func (s *Semaphore) Acquire(n int) {

	s.m.Lock()
	defer s.m.Unlock()
	if s.cur >= int64(n) {

		s.cur-=int64(n)
		s.take(n)
		return

	}


	//如果没有返回 这里依旧还有锁
	c:=make(chan struct{})
	//获取下一个
	if s.w == nil{
		s.w =&waiters{value: int64(n),c:c}
	}else{
		//如果存在数据
		s.w.next = &waiters{value: int64(n),c:c}
	}

	//阻塞住
	fmt.Println("阻塞了:",n)
	//睡眠之前 尝试通知其他waiter
	s.notifyWaiters()
	s.m.Unlock()  //睡眠之前先解锁 把机会让给别人


	<-c
	s.take(n)

	//return之前需要先加锁 defer还持有锁
	s.m.Lock()
	return
}

func (s *Semaphore) Release(n int) {

	if atomic.LoadInt64(&s.cur)+int64(n) > s.max {
		panic("N is too large")
	}
	s.m.Lock()
	s.cur +=int64(n)
	s.send(n)
	//唤醒阻塞队列中的数据
	s.notifyWaiters()
	s.m.Unlock()
}

func(s *Semaphore) notifyWaiters(){

		if s.w != nil {
			if s.cur < s.w.value {
				return //如果cur不足够获取，则直接跳过，避免饥饿
			}

			close(s.w.c) //唤醒
			s.cur -= int64(s.w.value)
			fmt.Println("唤醒了:", s.w.value)
			s.remove() //移除队列头部的waiter
		}

}


func(s *Semaphore) remove(){
	next:=s.w.next
	if next!=nil {
		s.w = next.(*waiters)
	}

	s.w = nil

}



func main(){
	s:=NewSemaphore(10)
	go func(){
		s.Acquire(8)
		time.Sleep(time.Second)
		s.Release(8)
	}()


	go func(){
		s.Acquire(9)
		time.Sleep(time.Second)
		s.Release(9)
	}()



	time.Sleep(30 * time.Second)
}