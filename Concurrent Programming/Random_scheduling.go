package main

import (
	"fmt"
	"sync"
	// "time"
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

func Routine() {
    wg := sync.WaitGroup{}
    wg.Add(10)
    for i := 0; i < 10; i++ {
        go func(n int) {
            defer wg.Done()
            fmt.Println(n)
        }(i)
    }
    wg.Wait()
}

func solve() {
    Routine()
}

func main() {
    t := 1
    for ; t > 0; t-- {
        solve()
    }
}
