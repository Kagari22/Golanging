package main

import (
	"context"
	"fmt"
	"sync"
)

func solve() {
    wg := sync.WaitGroup{}
    ctx1 := context.WithValue(context.Background(), "title", "Pig1")
    ctx2 := context.WithValue(ctx1, "title2", "Pig2")
    ctx3 := context.WithValue(ctx2, "title3", "Pig3")

    wg.Add(1)
    go func(c context.Context) {
        defer wg.Done()
        if v := c.Value("title"); v != nil {
            fmt.Println("Found value:", v)
            return
        }
        fmt.Println("Key not found:", "title")
    }(ctx3)

    wg.Wait()
}

func main() {
	solve()
}
