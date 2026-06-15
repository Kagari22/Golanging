package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var data []byte
var idx int

func nextInt() int {
    sign := 1
    for idx < len(data) && (data[idx] == ' ' || data[idx] == '\n' || data[idx] == '\r' || data[idx] == '\t') {
        idx++
    }
    if idx < len(data) && data[idx] == '-' {
        sign = -1
        idx++
    }
    num := 0
    for idx < len(data) && data[idx] >= '0' && data[idx] <= '9' {
        num = num*10 + int(data[idx]-'0')
        idx++
    }
    return sign * num
}

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
