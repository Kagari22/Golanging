package main

import (
	"fmt"
	"sync"
)

func solve() {
    wg := sync.WaitGroup{}
    var mp sync.Map
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            mp.Store(n, fmt.Sprintf("value: %d", n))
        }(i)
    }

    wg.Wait()

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            fmt.Println(mp.Load(n))
        }(i)
    }

    wg.Wait()

    mp.Range(func(key, value any) bool {
        fmt.Println(key, value)
        return true
    })

    mp.Delete(4)
}

func main() {
	solve()
}
