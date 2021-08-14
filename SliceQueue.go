package main

import "sync"

type SliceQueue struct{
	data []interface{}
	mu sync.Mutex
}

func NewSliceQueue(n int) (q *SliceQueue){
	return &SliceQueue{data:make([]interface{},0,n)}
}

//Enqueue 把值放在队列尾部
func (q *SliceQueue) Enqueue(v interface{}){
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = append(q.data,v)
}


//Dequeue 取出队头并返回
func (q *SliceQueue)Dequeue() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.data) ==0{
		return nil
	}
	v :=q.data[0]
	q.data = q.data[1:]
	return v
}