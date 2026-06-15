package main

import (
	"fmt"
	"sync"
	// "math/rand/v2"
	"time"
)

func solve() {
    wg := sync.WaitGroup{}
    ch := make(chan struct{})

    wg.Add(1)
    go func() {
        defer wg.Done()
        time.Sleep(2 * time.Second)
        fmt.Println("发出信号, close(ch)")
        close(ch)  
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            select {
            case <-ch:
                fmt.Println("收到信号, <-ch")
                return
            default:
            }
            fmt.Println("业务逻辑处理中...")
            time.Sleep(300 * time.Millisecond)
        }
    }()

    wg.Wait()
}

func main() {
    solve()
}
