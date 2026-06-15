package main

import (
	// "fmt"
)

func Channel() {
    // fmt.Println("Before select")
    // // 空 select 阻塞
    // select {}
    // fmt.Println("After select")

    // // nil channel 无法读写
    // var ch chan int 

    // fmt.Println("Before select")
    // select {
    // // 读写无法执行被阻塞
    // case <- ch:
    // case ch <- 1919810:
    // }
    // fmt.Println("After select")
}

func solve() {
    Channel()
}

func main() {
    t := 1
    for ; t > 0; t-- {
        solve()
    }
}
