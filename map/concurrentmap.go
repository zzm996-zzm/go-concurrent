package main

import (
 "fmt"
 "unsafe"
)

func main() {
 var expunged = unsafe.Pointer(new(interface{}))
 fmt.Println(expunged)
}
