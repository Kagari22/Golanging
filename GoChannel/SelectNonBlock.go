package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func Channel() {
    people := 10
    max := 10
    ans := rand.Intn(max)
    fmt.Println("The answer is: ", ans)
    fmt.Println("---------------------------")
    
    ch := make(chan int, people)
    wg := sync.WaitGroup{}
    wg.Add(people)
    for i := 0; i < people; i++ {
        go func() {
            defer wg.Done()
            res := rand.Intn(max)
            fmt.Println("Someone guessed: ", res)
            if res == ans {
                ch <- res
            }
        }()
    }
    wg.Wait()

    fmt.Println("---------------------------")
    // 非阻塞收发
    select {
    case res := <-ch:
        fmt.Println("Someone hit the answer: ", res)
    default:
        fmt.Println("No one hit the answer!")
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
