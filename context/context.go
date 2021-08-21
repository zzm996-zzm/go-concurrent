package main

import (
	"context"
	"fmt"
	"time"
)


func main(){
	base:=context.Background()
	cv:=context.WithValue(base,1,2)
	cl1,_:=context.WithCancel(cv)
	cl2,_:=context.WithCancel(cl1)
	ct,ct1:=context.WithTimeout(cl2,time.Second*30)
	cv2:=context.WithValue(ct,2,3)
	cv3:=context.WithValue(ct,3,4)


	time.Sleep(100)
	//第一个go
	go func(ctx context.Context){
		select{
		 case <-ctx.Done():
		 	fmt.Println("cv被取消了")
		}
	}(cv)

	time.Sleep(100)
	//第二个go
	go func(ctx context.Context){
		select{
		case <-ctx.Done():
			fmt.Println("cl1被取消了")
		}
	}(cl1)


	time.Sleep(100)
	//第三个go
	go func(ctx context.Context){
		select{
		case <-ctx.Done():
			fmt.Println("cl2被取消了")
		}
	}(cl2)


	time.Sleep(100)
	//第四个go
	go func(ctx context.Context){
		select{
		case <-ctx.Done():
			fmt.Println("ct被取消了")
		}
	}(ct)

	time.Sleep(100)
	//第五个go
	go func(ctx context.Context){
		select{
		case <-ctx.Done():
			fmt.Println("cv2被取消了")
		}
	}(cv2)

	time.Sleep(100)
	//第六个go
	go func(ctx context.Context){
		select{
		case <-ctx.Done():
			fmt.Println("cv3被取消了")
		}
	}(cv3)





	fmt.Println("执行ct1 cancle")
	time.Sleep(time.Second * 1)
	ct1()
	fmt.Println(cv3.Err())

}