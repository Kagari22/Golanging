package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Result struct {
	Index int
	Data  string
}

func AllWithChannel() {
	const Qnum = 5

	ch := make(chan Result, Qnum)

	for i := 0; i < Qnum; i++ {
		go func(i int) {
			delay := rand.Intn(1000)
			time.Sleep(time.Duration(delay) * time.Millisecond)

			ch <- Result{
				Index: i,
				Data:  fmt.Sprintf("data_from_%d", i),
			}
		}(i)
	}

	// 收集所有结果
	for i := 0; i < Qnum; i++ {
		r := <-ch
		fmt.Println("Got:", r)
	}
}

func All() {
    const Qnum = 5

	var wg sync.WaitGroup
	results := make([]Result, Qnum)

	wg.Add(Qnum)

	for i := 0; i < Qnum; i++ {
		go func(i int) {
			defer wg.Done()
			// 模拟不同耗时
			delay := rand.Intn(1000)
			fmt.Println("Worker", i, "start, need", delay, "ms")
			time.Sleep(time.Duration(delay) * time.Millisecond)
			// 写结果
			results[i] = Result{
				Index: i,
				Data:  fmt.Sprintf("data_from_%d", i),
			}
			fmt.Println("Worker", i, "done")
		}(i)
	}

	// 等所有 goroutine 完成
	wg.Wait()

	fmt.Println("\n✅ All results collected:")
	for _, r := range results {
		fmt.Println(r)
	}
}

func solve() {
    All()
}

func main() {
    t := 1
    for ; t > 0; t-- {
        solve()
    }
}
