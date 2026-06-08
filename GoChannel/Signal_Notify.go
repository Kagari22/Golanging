package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func solve() {
    go func() {
        for {
            fmt.Println(time.Now().Format("15.04.05.000"))
            time.Sleep(300 * time.Millisecond)
        }
    }()

    // select{}
    chInt := make(chan os.Signal, 1)
    signal.Notify(chInt, os.Interrupt)
    defer signal.Stop(chInt)

    <-chInt
    fmt.Println("os signal interrupt")
}

func main() {
    solve()
}
