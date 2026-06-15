package main

import (
	"fmt"
	"time"
)

func solve() {
    t := time.NewTimer(time.Second)
    fmt.Println("Set the timer. \ttime is ", time.Now().String())
    now := <-t.C
    fmt.Println("The time is up, time is ", now.String())
}

func main() {
    t := 1
    for ; t > 0; t-- {
        solve()
    }
}
