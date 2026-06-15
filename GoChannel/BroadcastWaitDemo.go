package main

import (
	"fmt"
	"sync"
	"time"
)

func solve() {
	wg := sync.WaitGroup{}
    datacap := 1024 * 1024
    var data []int
    cond := sync.NewCond(&sync.Mutex{})
    
    for i := 0; i < 8; i++ {
        wg.Add(1)
        go func(c *sync.Cond) {
            defer wg.Done()
            c.L.Lock()
            for len(data) < datacap {
                c.Wait()
            }
            fmt.Println("listen", len(data), time.Now())
            c.L.Unlock()
        }(cond)
    }

    wg.Add(1)   
    go func(c *sync.Cond) {
        defer wg.Done()
        c.L.Lock()
        defer c.L.Unlock()
        for i := 0; i < datacap; i++ {
            data = append(data, i * i)
        }
        fmt.Println("Broadcast")
        c.Broadcast()
    }(cond)

    wg.Wait()
}

func main() {
	solve()
}
