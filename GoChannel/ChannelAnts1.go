package main

import (
	"fmt"
	"runtime"
	"time"
)

// 用Channel来控制并发规模
func Channel() {
    go func() {
        for {
            fmt.Println("Num: ", runtime.NumGoroutine())
            time.Sleep(500 * time.Millisecond)
        }
    }()

    const size = 114
    ch := make(chan struct{}, size)
    
    for {
        ch <- struct{}{}
        go func() {
            time.Sleep(time.Second)
            <-ch
        }()
    }
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
