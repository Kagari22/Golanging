package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
	"github.com/panjf2000/ants/v2"
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
    go func() {
        for {
            fmt.Println("Num: ", runtime.NumGoroutine())
            time.Sleep(100 * time.Millisecond)
        }
    }()

    size := 1024
    // ants 控制 goroutine 并发数量
    pool, err := ants.NewPool(size)
    if err != nil {
        log.Fatalln(err)
    }
    defer pool.Release()

    for {
        err := pool.Submit(func() {
            time.Sleep(100 * time.Second)
        })
        if err != nil {
            log.Fatalln(err)
        }
    }
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
