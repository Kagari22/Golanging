package main

import (
	"fmt"
	"sync"
)

func solve() {
    wg := sync.WaitGroup{}
    mu := sync.Mutex{}

    cnt := 0
    t := 1000
    wg.Add(t)

    for i := 0; i < t; i++ {
        go func() {
            defer wg.Done()
            for j := 0; j < 100; j++ {
                mu.Lock()
                cnt++;
                mu.Unlock()
            }
        }()
    }

    wg.Wait()
    fmt.Println("Cnt:", cnt)
}

func main() {
	solve()
}
