package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	NodeName string
	Addr  string
	Count int32
}

func  loadNewConfig() Config{
	return  Config{
		NodeName :"北京",
		Addr : "10.77.95.27",
		Count : rand.Int31(),
	}
}

func main(){
	var config atomic.Value
	config.Store(loadNewConfig())
	var cond = sync.NewCond(&sync.Mutex{})

	//设置新的config
	go func(){
			for{
				time.Sleep(time.Duration(5+rand.Int63n(5)) * time.Second)
				config.Store(loadNewConfig())
				//通知所有等待获取变更信号的携程
				cond.Broadcast()
			}
	}()

	go func(){
		for{
			cond.L.Lock()
			cond.Wait()
			c:=config.Load().(Config)
			fmt.Printf("new config : %+v\n",c)
			cond.L.Unlock()
		}
	}()

	select{}
}