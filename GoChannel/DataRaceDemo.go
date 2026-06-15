package main

import (
	"fmt"
	"sync"
)

// go run -race .\test.go
func solve() {
	wg := sync.WaitGroup{}
    cnt := 0
    gs := 1000
    wg.Add(gs)

    for i := 0; i < gs; i++ {
        go func() {
            defer wg.Done()
            for k := 0; k < 100; k++ {
                cnt++;
            }
        }()
    }

    wg.Wait()
    fmt.Println("Cnt:", cnt)
}

func main() {
	solve()
}
