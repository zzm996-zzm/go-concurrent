// package main

// import (
// 	"fmt"
// 	"sync/atomic"
// 	"time"
// )

// type _atomic struct{
// 	atomic.Value
// }

// type config struct {
// 	id int
// 	name string
// 	arr  []int
// 	m  map[string]chan int
// }

// type Compare struct{
// 	i int
// 	s  string
// 	c  chan int

// }

// func compare(){

// 	a := Compare{
// 		1,"zzm",make(chan int),
// 	}

// 	b := Compare{
// 		2,"zzm", make(chan int),
// 	}
// 	if a == b{
// 		fmt.Println("true")
// 	}
// }

// func test(){
// 	var c chan int = make(chan int)

// 	go func(){
// 		i:=0
// 		for i<1{
// 			i++
// 			c<-1
// 			time.Sleep(1000)
// 		}
// 	}()

// 	<-c
// 	close(c)
// }

// func main() {
// 	ch1 := make(chan int)
// 	ch2 := make(chan int)
// 	ch3 := make(chan int)
// 	ch4 := make(chan int)
// 	go func() {
// 		for {
// 			fmt.Println("I'm goroutine 1")
// 			ch2 <-1 //I'm done, you turn
// 			<-ch1
// 		}
// 	}()

// 	go func() {
// 		for {
// 			<-ch2
// 			fmt.Println("I'm goroutine 2")
// 			ch3 <-1
// 		}

// 	}()

// 	go func() {
// 		for {
// 			<-ch3
// 			fmt.Println("I'm goroutine 3")
// 			ch4 <-1
// 		}

// 	}()

// 	go func() {
// 		for {
// 			<-ch4
// 			fmt.Println("I'm goroutine 4")
// 			ch1 <-1
// 		}

// 	}()

// 	select {}
// }