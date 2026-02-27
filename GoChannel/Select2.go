package main

import (
	"fmt"
	"math/rand"
	"time"
)

func Channel() {
   ch := make(chan int)

    go func() {
        for {
            ch <- rand.Intn(100)
            time.Sleep(500 * time.Millisecond)
        }
    }()

    go func() {
        for {
            select {
            case v := <-ch: 
                fmt.Println("Received value: ", v)
            }
        }
    }()

    time.Sleep(5 * time.Second)
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
