// 统计目录的文件数量和大小
// go run main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size) / GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size) / MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size) / KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// dir：当前目录
// wg：等待所有目录遍历完成
// sizes：发送文件大小的Channel
func walkDir(dir string, wg *sync.WaitGroup, sizes chan<- int64) {
	// 当前目录遍历结束
	defer wg.Done()

	// 读取目录下所有文件和子目录
	entries, err := os.ReadDir(dir)
	if (err != nil) {
		return 
	}

	// 遍历当前目录
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		// 如果是子目录
		if entry.IsDir() {
			// 新增一个任务
			wg.Add(1)
			// 开启新的Goroutine递归遍历
			go walkDir(path, wg, sizes)
		} else {
			// 普通文件
			// 获取文件信息
			info, err := entry.Info()
			if err != nil {
				continue
			}
			// 将文件大小发送到Channel
			sizes <- info.Size()
		}
	}
}
	
func main() {
	// 获取命令行参数
	// go run main.go dir1 dir2 dir3
	roots := os.Args[1:]
	
	// 默认统计当前目录
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// 用于传递文件大小
	sizes := make(chan int64)

	// 等待所有walkDir结束
	var wg sync.WaitGroup

	// 为每个根目录启动一个Goroutine
	for _, root := range roots {
		info, err := os.Stat(root)
		if (err != nil) {
			continue
		}
		// 如果是目录
		if info.IsDir() {
			wg.Add(1)
			go walkDir(root, &wg, sizes)
		} else {
			// 如果直接传入的是文件
			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				info, err := os.Stat(file)
				if (err == nil) {
					sizes <- info.Size()
				}
			}(root)
		}
	}

	// 等待所有目录遍历完成
	go func() {
		// 等待所有walkDir退出
		wg.Wait()
		// 通知Collector没有新的数据了
		close(sizes)
	}()

	var (
		fileCount int64
		totalSize int64
	)

	done := make(chan struct{})

	// 单独启动一个统计协程
	go func() {
		// 不断读取Channel中的文件大小
		for size := range sizes {
			fileCount++
			totalSize += size
		}
		// 统计完成
		close(done)
	}()

	// 等待统计结束
	<-done

	cwd, _ := os.Getwd()
	fmt.Println("当前工作目录：", cwd)
	fmt.Printf("%d files %s\n", fileCount, formatSize(totalSize))
}
