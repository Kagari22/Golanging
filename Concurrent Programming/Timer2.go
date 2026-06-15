package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func solve() {
    ch := make(chan int)

    go func() {
        for {
            ch <- rand.IntN(10)
            time.Sleep(400 * time.Millisecond)
        }
    }()

    // 每局时间
    t := time.NewTimer(3 * time.Second)
    hint, miss := 0, 0
    // 统计结果，共玩5次
    for i := 0; i < 5; i++ {
    guess:
        for {
            select {
            case v := <-ch:
                fmt.Println("Guess value: ", v)
                if v == 4 {
                    fmt.Println("Bingo! Some one hint the answer.")
                    // 新游戏，重置定时器
                    t.Reset(3 * time.Second)
                    hint++
                    break guess
                }
            case <-t.C:
                fmt.Println("The time is up, no one hint.")
                miss++
                // 重新创建定时器
                t = time.NewTimer(3 * time.Second)
                break guess
            }
        }
    }
    fmt.Println("Game over! Hint ", hint, " Miss ", miss)
}

func main() {
    t := 1
    for ; t > 0; t-- {
        solve()
    }
}
