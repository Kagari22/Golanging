package llm

import (
	"context"
	"fmt"
	"strings"
)

type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) Generate(_ context.Context, messages []Message, _ []ToolSpec) (Response, error) {
	if len(messages) == 0 {
		return Response{Content: "Hello. Ask me to calculate or tell the time."}, nil
	}

	last := messages[len(messages) - 1] // 取最后一条消息
	// 如果最后一条是工具结果
	if last.Role == RoleTool {
		switch last.Name {
		case "calculator":
			return Response{Content: fmt.Sprintf("The result is %s.", last.Content)}, nil
		case "time":
			return Response{Content: fmt.Sprintf("Current time: %s", last.Content)}, nil
		default:
			return Response{Content: last.Content}, nil
		}
	}

	// 最后一条不是用户消息
	if last.Role != RoleUser {
		return Response{Content: "Ready for the next message."}, nil
	}

	// 处理用户消息
	text := strings.TrimSpace(last.Content)
	lower := strings.ToLower(text)

	if strings.HasPrefix(lower, "calc ") {
		return Response{
			ToolCall: &ToolCall{
				Name:  "calculator",
				Input: strings.TrimSpace(text[5:]),
			},
		}, nil
	}

	if strings.Contains(lower, "time") || strings.Contains(lower, "clock") || strings.Contains(text, "几点") || strings.Contains(text, "时间") || strings.Contains(text, "现在几点") {
		return Response{
			ToolCall: &ToolCall{
				Name:  "time",
				Input: "local",
			},
		}, nil
	}

	return Response{
		Content: "Mock mode is active. Try 'calc 12*(3+4)' or ask me for the current time. Set OPENAI_API_KEY, OPENAI_BASE_URL, and OPENAI_MODEL to use a real model.",
	}, nil
}
