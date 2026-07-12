package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go-agent/agent"
	"go-agent/llm"
	"go-agent/memory"
	"go-agent/tools"
)

/*
$env:OPENAI_API_KEY="你的API Key"
$env:OPENAI_BASE_URL="https://api.openai.com/v1"
$env:OPENAI_MODEL="gpt-4o-mini"
go run .
*/

func main() {
	fmt.Println("Go Agent")
	fmt.Println("Type a message. Use 'exit' to quit.")

	// 创建一个对话记忆仓库, 并把最多保留的消息数限制为 24 条
	store := memory.NewStore(24)
	// 创建一个工具列表, 并把 agent 目前能用的两个工具注册进去
	toolset := []tools.Tool{
		tools.NewCalculatorTool(), // 返回一个计算器工具
		tools.NewTimeTool(),       // 返回一个时间工具
	}

	client := buildClient()
	runner := agent.NewRunner(client, store, toolset, 4) // 创建 agent 的主控制器

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nYou> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if strings.EqualFold(input, "exit") {
			fmt.Println("Bye.")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
		answer, err := runner.Run(ctx, input) // 让 agent 处理用户输入
		cancel() // 任务完成就立刻取消
		if err != nil {
			fmt.Printf("Agent error: %v\n", err)
			continue
		}

		fmt.Printf("Agent> %s\n", answer)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Input error: %v\n", err)
	}
}

func buildClient() llm.Client {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	baseURL := strings.TrimSpace(os.Getenv("OPENAI_BASE_URL"))
	model := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))

	if apiKey == "" || baseURL == "" || model == "" {
		return llm.NewMockClient()
	}

	return llm.NewOpenAICompatibleClient(apiKey, baseURL, model)
}
