package main

import (
	"errors"
	"fmt"
	"sync"
)

type Queue interface {
	Add(n int) error
	Pop() (int, error)
	Close() error
}

type LimitQueue struct {
	queue  []int
	cond   *sync.Cond
	closed chan struct{}
}

func NewLimitQueue(cap uint32, L sync.Locker) *LimitQueue {
	queue := make([]int, 0, cap)
	return &LimitQueue{queue: queue, cond: sync.NewCond(L)}

}

func (q *LimitQueue) Add(n int) error {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if len(q.queue) == 0 || len(q.queue) < cap(q.queue) {
		q.queue = append(q.queue, n)
		//唤醒所有waiter 通知队列中有元素了
		q.cond.Broadcast()
		return nil
	}
	return errors.New("Add失败，有限队列已满")
}

func (q *LimitQueue) Pop() (int, error) {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for len(q.queue) == 0 {
		select {
		case _, ok := <-q.closed:
			if ok {
				return 0, errors.New("当前队列已经关闭了")
			}
		default:
			fmt.Println("阻塞")
			q.cond.Wait()
			fmt.Println("被唤醒了")
		}

	}

	n := q.queue[0]
	q.queue = q.queue[1:]
	return n, nil
}
func (q *LimitQueue) Close() error {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if len(q.queue) != 0 {
		return errors.New("队列还存在值，关闭失败")
	}
	close(q.closed)
	q.cond.Broadcast() //关闭时通知等待的goroutine，避免它们永远等待
	return nil
}

// func main() {
// 	q := NewLimitQueue(10, &sync.Mutex{})

// 	go func() {
// 		x, _ := q.Pop()
// 		fmt.Println("go0", x)
// 	}()

// 	go func() {
// 		x, _ := q.Pop()
// 		fmt.Println("go1", x)
// 	}()

// 	go func() {

// 		for i := 0; i < 2; i++ {
// 			q.Add(i)
// 			time.Sleep(1 * time.Second)
// 			fmt.Println("add ", i)

// 		}
// 	}()
// 	time.Sleep(10 * time.Second)
// }
