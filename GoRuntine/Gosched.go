package main

import (
	"fmt"
	"runtime"
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

func Routine() {
    runtime.GOMAXPROCS(1)
    wg := sync.WaitGroup{}
    wg.Add(1)
    max := 100

    go func() {
        defer wg.Done()
        for i := 1; i < max; i += 2 {
            fmt.Print(i, " ")
            // 主动让出当前 goroutine 的时间片，把自己放回可运行队列
            runtime.Gosched()
        }
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 2; i < max; i += 2 {
            fmt.Print(i, " ")
            // 主动让出当前 goroutine 的时间片，把自己放回可运行队列
            runtime.Gosched()
        }
    }()

    wg.Wait()

    // 输出结果出现连续两个偶数/奇数 -> runtime.Gosched() 只让出 CPU，但调度器可能再次调度同一个 goroutine，因此不能保证严格交替执行。
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
