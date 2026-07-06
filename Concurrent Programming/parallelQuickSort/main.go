package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 大于该长度才开启并发
const Threshold = 10000 

// 生成随机数组
func randomArray(n, l, r int) []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(r-l+1) + l
	}
	return a
}

// 划分
func partition(a []int, l, r int) int {
	pivot := a[r]
	i := l

	for j := l; j < r; j++ {
		if a[j] <= pivot {
			a[i], a[j] = a[j], a[i]
			i++
		}
	}

	a[i], a[r] = a[r], a[i]
	return i
}

// 普通快速排序
func quickSort(a []int, l, r int) {
	if l >= r {
		return
	}

	p := partition(a, l, r)

	quickSort(a, l, p - 1)
	quickSort(a, p + 1, r)
}

// 并发快速排序
func parallelQuickSort(a []int, l, r int) {
	if l >= r {
		return
	}

	p := partition(a, l, r)

	if r - l < Threshold {
		quickSort(a, l, p - 1)
		quickSort(a, p + 1, r)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		parallelQuickSort(a, l, p - 1)
	}()

	go func() {
		defer wg.Done()
		parallelQuickSort(a, p + 1, r)
	}()

	wg.Wait() 
}

func check(a []int) bool {
	for i := 1; i < len(a); i++ {
		if a[i] < a[i - 1] {
			return false
		}
	}
	return true
}

func solve() {
	n := 10000000

	a := randomArray(n, 0, 1000000000)

	start := time.Now()

	// quickSort(a, 0, len(a) - 1)
	parallelQuickSort(a, 0, len(a) - 1)

	elapsed := time.Since(start)

	fmt.Println("排序完成")
	fmt.Println("是否有序", check(a))
	fmt.Println("耗时:", elapsed)
}

func main() {
	solve()
}
