package main

import (
	"fmt"
	"reflect"
)

func main(){
	var ch1 = make(chan int,10)
	var ch2 = make(chan int,10)

	var cases = createCases(ch1,ch2)
	//执行10次select
	for i:=0;i<10;i++{
		chosen,recv,ok :=reflect.Select(cases)
		if recv.IsValid(){
			fmt.Println("recv:", cases[chosen].Dir, recv, ok)
		}else{
			fmt.Println("send:", cases[chosen].Dir, ok)
		}
	}
}

func createCases(chs ...chan int)[]reflect.SelectCase{
	var cases []reflect.SelectCase

	//创建recv case
	for _,ch:=range chs{
		cases = append(cases,reflect.SelectCase{
			Dir: reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}

	//send case
	for i,ch:=range chs{
		v:=i
		cases = append(cases,reflect.SelectCase{
			Dir: reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: reflect.ValueOf(v),
		})
	}

	return cases


}