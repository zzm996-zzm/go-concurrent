package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// type Once struct {
// 	done uint32
// }

// //正常人第一时间反应想到的once实现
// /**
//  ** 如果f执行很慢，后续调用Do的方法会发现 done也就被atomic修改成1了,所以不会执行。
//  ** 然后获取初始化资源的时候可能会得到nil。因为f还没有执行完
// **/
// func (o *Once) Do(f func()) {
// 	if !atomic.CompareAndSwapUint32(&o.done, 0, 1) {
// 		return
// 	}
// 	f()
// }

type Once struct {
	done uint32
	sync.Mutex
}

func (o *Once) Done() bool {
	return atomic.LoadUint32(&o.done) == 1
}

func (o *Once) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 0 {
		return o.doSlow(f)
	}

	return nil
}

func (o *Once) doSlow(f func() error) error {
	o.Lock()
	defer o.Unlock()
	var err error
	// 双检查
	if o.done == 0 {
		err = f()
		if err != nil {
			return err
		}
	}

	atomic.StoreUint32(&o.done, 1)
	return nil

}

// Once 是一个扩展的sync.Once类型，提供了一个Done方法
type OnceTest struct {
	sync.Once
}

// Done 返回此Once是否执行过
// 如果执行过则返回true
// 如果没有执行过或者正在执行，返回false
func (o *OnceTest) Done() bool {
	return atomic.LoadUint32((*uint32)(unsafe.Pointer(&o.Once))) == 1
}

func main() {
	var flag OnceTest
	fmt.Println(flag.Done()) //false

	flag.Do(func() {
		time.Sleep(time.Second)
	})

	fmt.Println(flag.Done()) //true
}
