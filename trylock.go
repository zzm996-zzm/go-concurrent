package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type TryMutex struct {
	c        chan struct{}
	tryCount uint32
	state    int32
}
func (try *TryMutex) GetCount() uint32{
	return try.tryCount

}

func (try *TryMutex) TryLock() bool {

	if atomic.LoadInt32(&try.state) == 0 {

		if atomic.CompareAndSwapInt32(&try.state, 0, 1) {
			try.c = make(chan struct{}, 1)
		}
	}

	select {
	case try.c <- struct{}{}:
		return true
	case <-time.After(time.Second * 1):

		atomic.AddUint32(&try.tryCount, 1) //try.tryCount 并非原子操作
		return false
	default:
		atomic.AddUint32(&try.tryCount, 1)
		return false
	}

}

func main() {
	try := &TryMutex{}
	for i := 0; i < 1000; i++ {
		go func(){
			x := try.TryLock()
			fmt.Println(x)
		}()
	}

	time.Sleep(2 * time.Second)
	fmt.Println( try.GetCount())
}
