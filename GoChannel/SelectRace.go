package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func SelectRace() {
    type Rows struct {
        Index int
    }
    const Qnum = 8

    // 用于通信的channel，数据，停止信号
    ch := make(chan Rows, 1)
    stopChs := [Qnum]chan struct {}{}
    for i := range stopChs {
        stopChs[i] = make(chan struct{})
    }

    wg := sync.WaitGroup{}
    // 模拟 querier 查询，每个查询持续不同的时间
    wg.Add(Qnum)

    for i := 0; i < Qnum; i++ {
        // 每一个 querier
        go func(i int) {
            defer wg.Done()
            // 模拟执行时间
            randD := rand.Intn(1000)
            fmt.Println("Querier ", i, " starts fetch data, need duration is ", randD, " ms.")
            // 查询结果的channel
            chRst := make(chan Rows, 1)

            // 执行查询工作
            go func() {
                // 模拟时长
                time.Sleep(time.Duration(randD) * time.Millisecond)
                chRst <- Rows {
                    Index: i,
                }
            } ()

            // 监听查询结果和停止信号 channel
            select {
            // 查询结果
            case rows := <-chRst:
                fmt.Println("Querier ", i, " gets result.")
                // 保证没有其他结果写入，才写入结果
                if len(ch) == 0 {
                    ch <- rows
                }
            // stop信号
            case <-stopChs[i]: 
                fmt.Println("Querier ", i, " is stopping.")
                return
            }
        }(i)
    }

    // 等待第一个查询结果的反馈
    wg.Add(1)
    go func() {
        defer wg.Done()
        // 等待ch中传递的结果
        select {    
        // 等待第一个查询结果
        case rows := <-ch:
            fmt.Println("Get first result from ", rows.Index, ". Stop other querier.")
            // 循环结构，全部通知 querier 结束
            for i := range stopChs {
                // 当前返回结果的 goroutine 不需要了，因为已经结束
                if i == rows.Index {
                    continue
                }
                stopChs[i] <- struct{}{}
            }
        // 计划一个超时时间
        case <-time.After(5 * time.Second):
            fmt.Println("All queriers timeout.")
            // 循环结构，全部通知 querier 结束
            for i := range stopChs {
                stopChs[i] <- struct{}{}
            }
        }
    }()

    wg.Wait()
}

func solve() {
    SelectRace()
}

func main() {
    t := 1
    for ; t > 0; t-- {
        solve()
    }
}
