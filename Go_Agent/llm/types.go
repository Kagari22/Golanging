package llm

import "context"

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleTool      = "tool"
)

type Message struct {
	Role    string
	Name    string
	Content string
}

// 工具说明
type ToolSpec struct {
	Name        string
	Description string
}

// 工具调用
type ToolCall struct {
	Name  string
	Input string
}

// 模型返回的结果
type Response struct {
	Content  string
	ToolCall *ToolCall
}

type Client interface {
	Generate(ctx context.Context, messages []Message, tools []ToolSpec) (Response, error)
}
