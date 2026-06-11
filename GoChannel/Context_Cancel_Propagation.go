package main

import (
	"context"
	"fmt"
	"sync"
)

func solve() {
    ctx1, cancel := context.WithCancel(context.Background())
    ctx2, _ := context.WithCancel(ctx1)
    ctx3, _ := context.WithCancel(ctx1)
    ctx4, _ := context.WithCancel(ctx2)

    wg := sync.WaitGroup{}
    wg.Add(4)

    go func() {
        defer wg.Done()
        <-ctx1.Done()
        fmt.Println("one cancel")
    }()

    go func() {
        defer wg.Done()
        <-ctx2.Done()
        fmt.Println("two cancel")
    }()

    go func() {
        defer wg.Done()
        <-ctx3.Done()
        fmt.Println("three cancel")
    }()

    go func() {
        defer wg.Done()
        <-ctx4.Done()
        fmt.Println("four cancel")
    }()

    cancel()
    wg.Wait()
}

func main() {
	solve()
}
