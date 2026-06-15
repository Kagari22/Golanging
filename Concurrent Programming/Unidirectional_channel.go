package main

import (
	"fmt"
	"sync"
)

func getEle(ch <- chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Println("The received value is: ", <-ch)
}   

func setEle(ch chan <- int, v int, wg *sync.WaitGroup) {
    defer wg.Done()
    ch <- v
    fmt.Println("The sent value is: ", v)
}

func Channel() {
    ch := make(chan int)
    wg := sync.WaitGroup{}

    wg.Add(2)
    go setEle(ch, 114514, &wg)
    go getEle(ch, &wg)

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
