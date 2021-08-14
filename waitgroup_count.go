package main

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type State struct {
	oCopy  noCopy
	state1 [3]uint32
}

//由于是64位机器 没办法测试32位机器 64位原子性
// func main() {

// 	s := &State{}
// 	var x uint32 = 1
// 	fmt.Printf("%p \n", &s.state1)
// 	fmt.Printf("%p \n", &x)
// 	fmt.Println(uintptr(unsafe.Pointer(&s.state1)) % 8)
// 	state := (*uint64)(unsafe.Pointer(&s.state1))
// 	atomic.AddUint64(state, uint64(10)<<32)
// 	fmt.Println(s.state1[0], s.state1[1])
// }

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		go func() {
			wg.Add(1) // 计数值加1
			//wg.Done() // 计数器减1
		}()

	}

	time.Sleep(time.Second * 2)
	go func() {
		if uintptr(unsafe.Pointer(&wg))%8 == 0 {
			state := (*uint64)(unsafe.Pointer(&wg))
			fmt.Println("64 :")
			fmt.Println(*state >> 32)
			fmt.Println(uint32(*state))
		} else {
			state := (*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(&wg)) + uintptr(4)))
			fmt.Println("32 :")
			fmt.Println(*state >> 32)
			fmt.Println(uint32(*state))

		}
	}()

	time.Sleep(time.Second * 10)
	wg.Add(-100)
	wg.Wait() // 主goroutine等待，有可能和第7行并发执行
}
