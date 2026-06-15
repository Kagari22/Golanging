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
				// 调用 Wait() 后内部会释放锁，并将该 goruntine 挂起并等待唤醒
				// 当其他 goroutine 调用 cond.Broadcast() 唤醒它时，Wait() 会重新尝试获取锁（阻塞直到获得锁），获取锁成功后，Wait() 返回。
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
