package main

import (
	"fmt"
	"time"
)

func Channel() {
    const size = 3
    ch := make(chan struct{}, size)

    for i := 0; ; i++ {
        fmt.Println("准备发送: ", i)
        ch <- struct{}{}
        fmt.Println("发送成功: ", i)
        go func(id int) {
            time.Sleep(2 * time.Second)
            fmt.Println("释放: ", i)
            <-ch
        }(i)
    }
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
