package main

import (
	"fmt"
	"time"
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
    ch := make(chan int)

    go func() {
        ch <- 114 + 514
    }()

    go func() {
        v := <-ch
        fmt.Println("Value: ", v)
    }()

    time.Sleep(time.Second)
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
