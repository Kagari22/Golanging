package main

import (
	"fmt"
	"time"
)

func f() int {
    fmt.Println("f() was run")
    return 2
}

func Channel() {
    var a[4] int
    var c1, c2, c3, c4 = make(chan int), make(chan int), make(chan int), make(chan int)
    var i1, i2 int

    go func() {
        c1 <- 10
    }()

    go func() {
        <-c2
    }()

    go func() {
        close(c3)
    }()

    go func() {
        c4 <- 30
    }()

    go func() {
        select {
        case i1 = <-c1:
            fmt.Println("Received ", i1, " from c1")
        case c2 <- i2:
            fmt.Println("Sent ", i2, " to c2")
        case i3, ok := <-c3:
            if ok {
                fmt.Println("Received ", i3, " from c3")
            } else {
                fmt.Println("c3 is closed")
            }
        case a[f()] = <- c4: 
            fmt.Println("Received ", a[f()], " from c4")
        default:
            fmt.Println("No communication")
        }
    }()

    time.Sleep(time.Second)
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
