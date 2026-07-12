# Go Agent

一个最小可运行的 Go Agent 示例，默认无需 API Key 就能运行。

## 功能

- 命令行多轮对话
- 内存中的会话历史
- `calculator` 工具
- `time` 工具
- 默认 `mock` 模式
- 可切换到 OpenAI 兼容接口

## 目录

```text
Go_Agent/
  main.go
  go.mod
  agent/
  llm/
  memory/
  tools/
```

## 运行

```powershell
cd C:\Users\pc\Desktop\Go_Agent
go run .
```

示例：

```text
You> calc 12*(3+4)
Agent> The result is 84.

You> 现在几点
Agent> Current time: 2026-07-10T...
```

## 使用真实模型

设置下面三个环境变量后，程序会改用兼容 `chat/completions` 的接口：

```powershell
$env:OPENAI_API_KEY="your_api_key"
$env:OPENAI_BASE_URL="https://api.openai.com/v1"
$env:OPENAI_MODEL="gpt-4o-mini"
go run .
```

当前这个最小版本通过提示词约定 `TOOL:<name>:<input>` 来触发工具调用，目的是先把 Agent 闭环跑通。

## 后续可以怎么扩展

- 把 mock client 换成更完整的 Responses API / function calling
- 增加文件搜索、HTTP 请求、数据库查询等工具
- 增加日志、重试、工具超时
- 把内存改成 SQLite 或 JSON 持久化
