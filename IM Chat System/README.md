# IM Chat System

这是一个用来学习 Go 的最小 IM 聊天系统后端，包含：

- 用户注册
- 用户登录
- 简单 JWT 鉴权
- WebSocket 长连接
- 在线用户管理
- 单聊文本消息
- 消息保存到本地 JSON 文件
- 历史消息 / 离线消息拉取
- 简单网页测试客户端

## 运行

```bash
go mod tidy
go run ./cmd/server
```

然后打开：

```text
http://localhost:8080
```

默认数据会保存到：

```text
data.json
```

## API

### 注册

```text
POST /api/register
```

```json
{
  "username": "alice",
  "password": "123456",
  "nickname": "Alice"
}
```

### 登录

```text
POST /api/login
```

```json
{
  "username": "alice",
  "password": "123456"
}
```

### 当前用户

```text
GET /api/me
Authorization: Bearer <token>
```

### 历史消息

```text
GET /api/messages?peer_id=2
Authorization: Bearer <token>
```

### WebSocket

```text
ws://localhost:8080/ws?token=<token>
```

发送消息格式：

```json
{
  "type": "chat",
  "to": 2,
  "content": "hello"
}
```
