package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)
const (
	mutexLocked = 1 << iota //加锁标识位置
	mutexWoken 				//唤醒标识位置
	mutexStarving 			//锁饥饿标识位置
	mutexWaiterShift = iota //标识waiter的起始bit位置
)

type TryMuext struct{
	sync.Mutex
}

func (m *TryMuext) TryLock() bool{
	//如果能拿到锁
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)),0,mutexLocked){
		return true
	}

	//如果处于唤醒，加锁或者饥饿状态，这次请求就不参与竞争了。返回false
	//如果是饥饿状态那么又mutex去管理锁的分配，如果woken不为0表示有等待的获取锁的goroutine 不强制获取
	old :=atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexStarving|mutexWoken) != 0{
		return false
	}

	//尝试在竞争的状态下请求锁
	new:= old | mutexLocked //先把加锁位置设置为1
	//返回true 并且已经把lock加锁

	return  atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)),old,new)
}

func main(){
	var mu TryMuext
	go func(){
		mu.Lock()
		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
		mu.Unlock()
	}()

	time.Sleep(time.Second)
	ok:=mu.TryLock()
	if ok{
		fmt.Println("got the lock")
		mu.Unlock()
		return
	}

}