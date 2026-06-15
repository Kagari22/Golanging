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
			// Wait() 返回后不会直接执行 fmt.Println，而是先回到循环头部重新评估条件
			// 正是使用 for 而非 if 的原因，确保在每次唤醒后都重新检查共享状态。
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
