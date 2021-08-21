package main

import (
	"bytes"
	"sync"
)
type Copy interface {
   Test()
}


var buffers = sync.Pool{
	New :func() interface{}{
		return bytes.NewBuffer(make([]byte,1024,1024)) //1KB
	},
}

func init(){
	b1 := bytes.NewBuffer(make([]byte,1024,1024)) //1KB
	b2 := bytes.NewBuffer(make([]byte,2048,2048)) //2KB
	buffers.Put(b1)
	buffers.Put(b2)
}

func main(){

   //fmt.Println(buffers.Get().(*bytes.Buffer).Len())
	//fmt.Println(buffers.Get().(*bytes.Buffer).Len())
}
