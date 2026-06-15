package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func solve() {
	config := make(map[string]string)
    once := sync.Once{}

    loadConfig := func() {
        once.Do(func() {
            config = map[string]string{
                "varInt": fmt.Sprintf("%d", rand.Int31()),
            }
            fmt.Println("config loaded")
        }) 
    }

    workers := 10
    wg := sync.WaitGroup{}
    wg.Add(workers)

    for i := 0; i < workers; i++ {
        go func() {
            defer wg.Done()
            loadConfig()
            _ = config
        }()
    }

    wg.Wait()
}

func main() {
	solve()
}
