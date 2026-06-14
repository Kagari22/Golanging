package main

import (
	"fmt"
	"sync"
	"time"
)

var rw sync.RWMutex
var score = 90

func readScore(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    rw.RLock()
    fmt.Println("学生", id, "查到成绩:", score)
    rw.RUnlock()
}

func writeScore(newScore int) {
    rw.Lock()
    fmt.Println("老师正在改成绩...")
    score = newScore
    rw.Unlock()
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go readScore(i, &wg)
        time.Sleep(300 * time.Millisecond)
    }
    wg.Wait()

    writeScore(95)

    for i := 6; i <= 10; i++ {
        wg.Add(1)
        go readScore(i, &wg)
        time.Sleep(300 * time.Millisecond)
    }
    wg.Wait()
}
