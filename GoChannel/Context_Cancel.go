package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

func solve() {
    ctx, cancel := context.WithCancel(context.Background())
    var wg sync.WaitGroup

    for i := 0; i < 4; i++ {
        wg.Add(1)
        go func(c context.Context, n int) {
            defer wg.Done()
            for {
                select {
                case <-c.Done():
                    fmt.Println(strings.Repeat(" ", n), n)
                    return
                default:
                }
                fmt.Println(strings.Repeat(" ", n), n)
                time.Sleep(300 * time.Millisecond)
            }
        }(ctx, i)
    }

    time.Sleep(3 * time.Second)
    cancel()
    wg.Wait()
}

func main() {
	solve()
}
