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

// 异步Channel
// 接收和发送方不用等待双方就绪即可成功
func Channel() {
    // 设置缓冲区大小来保存发送的数据
    ch := make(chan int, 5)
    wg := sync.WaitGroup{}

    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 0; i < 5; i++ {
            ch <- i
            fmt.Println("Send ", i, "\tNow:", time.Now().Format("15:04:05.999999999"))
            time.Sleep(time.Second)
        }
        close(ch)
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        for e := range ch {
            fmt.Println("Received ", e, "\tNow:", time.Now().Format("15:04:05.999999999"))
            time.Sleep(4 * time.Second)
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
