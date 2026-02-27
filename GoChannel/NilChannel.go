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
            ch <- rand.Intn(10)
            time.Sleep(500 * time.Millisecond)
        }
    }()

    go func() {
        sum := 0
        t := time.After(5 * time.Second)
        for {
            select {
            case v := <-ch:
                fmt.Println("Received value: ", v)
                sum += v
            case <-t:
                // 将 channel 设置为 nil, 不再读写
                ch = nil
                fmt.Println("Ch was set nil, sum is: ", sum)
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
