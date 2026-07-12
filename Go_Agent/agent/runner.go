package agent

import (
	"context"
	"fmt"

	"go-agent/llm"
	"go-agent/memory"
	"go-agent/tools"
)

type Runner struct {
	client   llm.Client
	memory   *memory.Store
	tools    map[string]tools.Tool
	maxSteps int
}

// 把模型、记忆、工具集和最大步数组装成一个可运行的 agent 控制器
func NewRunner(client llm.Client, memory *memory.Store, toolset []tools.Tool, maxSteps int) *Runner {
	toolMap := make(map[string]tools.Tool, len(toolset))
	for _, tool := range toolset {
		toolMap[tool.Name()] = tool
	}

	return &Runner{
		client:   client,
		memory:   memory,
		tools:    toolMap,
		maxSteps: maxSteps,
	}
}

func (r *Runner) Run(ctx context.Context, userInput string) (string, error) {
	r.memory.Add(llm.Message{Role: llm.RoleUser, Content: userInput}) // 把用户消息写进记忆

	for step := 0; step < r.maxSteps; step++ {
		// 把当前聊天上下文和工具列表交给模型, 让模型决定是直接回答, 还是先调用某个工具
		response, err := r.client.Generate(ctx, r.memory.Messages(), availableTools(r.tools))
		if err != nil {
			return "", err
		}

		// 模型已经有最终答案, 不需要再调工具
		if response.ToolCall == nil {
			r.memory.Add(llm.Message{Role: llm.RoleAssistant, Content: response.Content})
			return response.Content, nil
		}

		// 检查工具是否存在
		tool, ok := r.tools[response.ToolCall.Name]
		if !ok {
			return "", fmt.Errorf("unknown tool: %s", response.ToolCall.Name)
		}

		// 保存 agent 的调用工具决策
		r.memory.Add(llm.Message{
			Role:    llm.RoleAssistant,
			Content: fmt.Sprintf("Calling tool %q with input %q", response.ToolCall.Name, response.ToolCall.Input),
		})

		// 执行工具
		result, err := tool.Run(ctx, response.ToolCall.Input)
		if err != nil {
			return "", fmt.Errorf("tool %s failed: %w", response.ToolCall.Name, err)
		}

		// 保存工具执行结果
		r.memory.Add(llm.Message{
			Role:    llm.RoleTool,
			Name:    response.ToolCall.Name,
			Content: result,
		})
	}

	return "", fmt.Errorf("agent stopped after %d steps", r.maxSteps)
}

// 把程序内部工具集合转换成模型可理解的工具列表
func availableTools(toolMap map[string]tools.Tool) []llm.ToolSpec {
	specs := make([]llm.ToolSpec, 0, len(toolMap))
	for _, tool := range toolMap {
		specs = append(specs, llm.ToolSpec{
			Name:        tool.Name(),
			Description: tool.Description(),
		})
	}
	return specs
}
