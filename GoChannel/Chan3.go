package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func Channel() {
    ch := make(chan int)
    wg := sync.WaitGroup{}

    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 0; i < 5; i++ {
            ch <- rand.Intn(10)
        }
        close(ch)
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        for e := range ch {
            fmt.Println("Value: ", e)
        }
    }()

    wg.Wait()
}

func solve() {
    Channel()
}

func main() {
    t := 1
    for ; t > 0; t-- {
        solve()
    }
}
