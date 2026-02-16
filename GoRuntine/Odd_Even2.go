package main

import (
	"fmt"
	"sync"
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

// 用 sync.WaitGroup 让 main 等待两个 goroutine 结束
func Routine() {
    var wg sync.WaitGroup
    wg.Add(2)

    oddf := func() {
        defer wg.Done()
        for i := 1; i <= 10; i += 2 {
            fmt.Println(i)
            time.Sleep(10 * time.Millisecond)
        }
    }
    evenf := func() {
        defer wg.Done()
        for i := 2; i <= 10; i += 2 {
            fmt.Println(i)
            time.Sleep(10 * time.Millisecond)
        }
    }
    
    go oddf()
    go evenf()

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
