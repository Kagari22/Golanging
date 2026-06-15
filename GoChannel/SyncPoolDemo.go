package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
)

func solve() {
	var cnt int32 = 0

    eleNewer := func() any {
        atomic.AddInt32(&cnt, 1)
        return new(bytes.Buffer)
    }

    pool := sync.Pool {
        New: eleNewer,
    }

    wg := sync.WaitGroup{}
    wg.Add(1000000)
    for i := 0; i < 1000000; i++ {
        go func() {
            defer wg.Done()
            buffer := pool.Get().(*bytes.Buffer)
            defer pool.Put(buffer)
            _ = buffer.String()
        }()
    }

    wg.Wait()
    fmt.Println("elements number is :", cnt)
}

func main() {
	solve()
}
