package main

import (
	"fmt"
	"sync"
	"unsafe"
)

type A int

type P *A
// 一个组合的并发原语
type MuOnce struct {
	sync.RWMutex
	sync.Once
}

func (m *MuOnce)Done() bool {
	p := (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(m)) + unsafe.Sizeof(sync.RWMutex{})))
	return *p == 1
}
func main() {
	m := MuOnce{}
	//p := (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&m))+unsafe.Sizeof(sync.RWMutex{})))
	fmt.Println(m.Done())
	m.Do(func(){
		fmt.Println("do")
	})
	fmt.Println(m.Done())

	//fmt.Println(*p)
}
