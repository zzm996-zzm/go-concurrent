package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type RecursiveMutex struct {
	sync.Mutex
	owner     int64 //当前持有锁的goroutine id
	recursion int32 //这个goroutine重入的次数
}

func (m *RecursiveMutex) Lock() {
	gid := GoID()
	//如果当前持有锁的goroutine就是这次调用的goroutine，说明是重入
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	//获得锁的goroutine第一次调用，记录下他的goroutine id
	atomic.StoreInt64(&m.owner, gid)

	m.recursion = 1

}

func (m *RecursiveMutex) Unlock() {
	gid := GoID()
	//非持有锁的goroutine 尝试释放锁，错误的使用
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	//调用次数减少1
	m.recursion--
	if m.recursion != 0 {
		return
	}

	// 此goroutine最后一次调用，需要释放锁
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()

}

func GoID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// 得到id字符串
	//Fields 去掉字符串中所有空格
	//TrimPrefix 去掉字符串中第二个参数
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return int64(id)
}

// Token方式的递归锁
type TokenRecursiveMutex struct {
	sync.Mutex
	token     int64
	recursion int32
}

// 请求锁，需要传入token
func (m *TokenRecursiveMutex) Lock(token int64) {
	if atomic.LoadInt64(&m.token) == token { //如果传入的token和持有锁的token一致，说明是递归调用
		m.recursion++
		return
	}
	m.Mutex.Lock() // 传入的token不一致，说明不是递归调用
	// 抢到锁之后记录这个token
	atomic.StoreInt64(&m.token, token)
	m.recursion = 1
}

// 释放锁
func (m *TokenRecursiveMutex) Unlock(token int64) {
	if atomic.LoadInt64(&m.token) != token { // 释放其它token持有的锁
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.token, token))
	}
	m.recursion--         // 当前持有这个锁的token释放锁
	if m.recursion != 0 { // 还没有回退到最初的递归调用
		return
	}
	atomic.StoreInt64(&m.token, 0) // 没有递归调用了，释放锁
	m.Mutex.Unlock()
}
