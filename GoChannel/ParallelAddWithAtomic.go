package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func solve() {
    cnt := atomic.Int32{}
    wg := sync.WaitGroup{}

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for i := 0; i < 100; i++ {
                cnt.Add(1)
            }
        }()
    }

    wg.Wait()
    fmt.Println("Cnt:", cnt.Load())
}

func main() {
	solve()
}
