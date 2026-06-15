package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

func solve() {
	var loadConfig = func() map[string]string {
		return map[string]string{
			"title":   "!?猪猪?!",
			"varConf": fmt.Sprintf("%d", rand.Int31()),
		}
	}

	var config atomic.Value

	go func() {
		for {
			config.Store(loadConfig())
			fmt.Println("last config was loaded", time.Now().Format("15:04:05.99999999"))
			time.Sleep(time.Second)
		}
	}()

	for {
		go func() {
			c := config.Load()
			fmt.Println(c, time.Now().Format("15:04:05.99999999"))
		}()
        time.Sleep(400 * time.Millisecond)
	}
}

func main() {
	solve()
}
