package main

import (
	"fmt"
	// "math/rand/v2"
	"time"
)

func solve() {
    ticker := time.NewTicker(time.Second)

    timer := time.After(5 * time.Second)

loop:
    for now := range ticker.C {
        fmt.Println("Now is", now.String())
        fmt.Println("http.Get(\"/ping\")")
        select {
        case <-timer:
            ticker.Stop()
            break loop
        default:
        }
    }
}

func main() {
    solve()
}
