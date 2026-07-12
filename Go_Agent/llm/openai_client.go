package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OpenAICompatibleClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewOpenAICompatibleClient(apiKey, baseURL, model string) *OpenAICompatibleClient {
	return &OpenAICompatibleClient{
		apiKey:     apiKey,
		baseURL:    strings.TrimRight(baseURL, "/"),
		model:      model,
		httpClient: &http.Client{},
	}
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *OpenAICompatibleClient) Generate(ctx context.Context, messages []Message, _ []ToolSpec) (Response, error) {
	// 组装请求体
	payload := chatRequest{
		Model:    c.model,
		Messages: make([]chatMessage, 0, len(messages) + 1),
	}

	// 加 system prompt
	payload.Messages = append(payload.Messages, chatMessage{
		Role: RoleSystem,
		Content: "You are a minimal Go agent. If a tool is needed, reply exactly with TOOL:<name>:<input>. " +
			"Available tools are calculator and time. Otherwise answer normally.",
	})

	// 把历史消息转成 API 消息
	for _, msg := range messages {
		payload.Messages = append(payload.Messages, chatMessage{
			Role:    normalizeRole(msg.Role),
			Content: flattenMessage(msg),
		})
	}

	// 把请求转成 JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return Response{}, err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL + "/chat/completions", bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	// 设置请求头
	req.Header.Set("Authorization", "Bearer " + c.apiKey) // 用 API Key 鉴权
	req.Header.Set("Content-Type", "application/json") // 告诉服务端这是 JSON

	// 发请求
	resp, err := c.httpClient.Do(req) 
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	// 处理 HTTP 错误
	if resp.StatusCode >= 300 {
		return Response{}, fmt.Errorf("model request failed: %s", strings.TrimSpace(string(data)))
	}

	// 解析 JSON
	var decoded chatResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		return Response{}, err
	}

	// 检查是否有内容
	if len(decoded.Choices) == 0 {
		return Response{}, fmt.Errorf("empty model response")
	}

	// 取出模型文本
	content := strings.TrimSpace(decoded.Choices[0].Message.Content)
	// 判断是不是工具调用
	if strings.HasPrefix(content, "TOOL:") {
		name, input, ok := parseToolDirective(content)
		if !ok {
			return Response{}, fmt.Errorf("invalid tool directive: %s", content)
		}
		return Response{
			ToolCall: &ToolCall{Name: name, Input: input},
		}, nil
	}

	return Response{Content: content}, nil
}

// 把 tool 角色转换成 user
func normalizeRole(role string) string {
	switch role {
	case RoleTool:
		return RoleUser
	default:
		return role
	}
}

// 把内部消息拍平成文本
func flattenMessage(msg Message) string {
	if msg.Role == RoleTool {
		return fmt.Sprintf("Tool %s returned: %s", msg.Name, msg.Content)
	}
	return msg.Content
}

// 把模型返回的工具调用文本解析成"工具名 + 工具输入"
func parseToolDirective(content string) (string, string, bool) {
	parts := strings.SplitN(content, ":", 3)
	if len(parts) != 3 {
		return "", "", false
	}
	name := strings.TrimSpace(parts[1])
	input := strings.TrimSpace(parts[2])
	if name == "" {
		return "", "", false
	}
	return name, input, true
}
