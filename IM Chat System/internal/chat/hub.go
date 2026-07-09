package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"IM_Chat_System/internal/auth"
	"IM_Chat_System/internal/store"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 70 * time.Second
	pingPeriod = 30 * time.Second
)

type Hub struct {
	mu      sync.RWMutex
	clients map[int64]*Client // 在线用户表, key 是用户 ID, value 是对应连接
	store   *store.Store // 消息和用户存储
	secret  string // 用来解析 token 的密钥
}

// 在线客户端连接
type Client struct {
	userID int64 // 这个连接属于哪个用户
	conn   *websocket.Conn // WebSocket 连接本身
	send   chan any // 发送队列, 往这里塞消息, writePump 会负责真正写到网络
	hub    *Hub // 指向所属的 Hub
}

// 前端发给服务端的消息格式
type IncomingMessage struct {
	Type    string `json:"type"`
	To      int64  `json:"to"`
	Content string `json:"content"`
}

// 服务端发回去的格式
type OutgoingMessage struct {
	Type    string        `json:"type"`
	Message store.Message `json:"message,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// 创建一个新的聊天中心
func NewHub(s *store.Store, secret string) *Hub {
	return &Hub{
		clients: make(map[int64]*Client),
		store:   s,
		secret:  secret,
	}
}

// 把普通 HTTP 请求升级成 WebSocket 连接
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claims, err := auth.ParseToken(h.secret, token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// 升级 HTTP 为 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade:", err)
		return
	}

	client := &Client{
		userID: claims.UserID,
		conn:   conn,
		send:   make(chan any, 16),
		hub:    h,
	}

	h.register(client)
	go client.writePump() // 读客户端发来的消息
	go client.readPump() // 把服务端消息写给客户端
}

// 把用户放进在线列表
func (h *Hub) register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if old := h.clients[client.userID]; old != nil {
		old.conn.Close()
	}
	h.clients[client.userID] = client
	log.Printf("user %d online\n", client.userID)
}

func (h *Hub) unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.userID] == client {
		delete(h.clients, client.userID)
		close(client.send)
		log.Printf("user %d offline\n", client.userID)
	}
}

// 尝试把消息投递给某个在线用户
func (h *Hub) deliver(userID int64, payload any) bool {
	h.mu.RLock()
	client := h.clients[userID]
	h.mu.RUnlock()

	if client == nil {
		return false
	}

	select {
	case client.send <- payload:
		return true
	default:
		client.conn.Close()
		return false
	}
}

// 持续读取客户端发来的消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(4096) // 最多读取 4096 字节
	c.conn.SetReadDeadline(time.Now().Add(pongWait)) // 设置超时时间
	// 注册一个函数, 当收到客户端发送的 Pong 帧时, 就自动执行这个函数
	c.conn.SetPongHandler(func(string) error { 
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var input IncomingMessage
		// 从 WebSocket 读取一条客户端 JSON 消息, 成功就存进 input, 失败就判断是否异常断开并结束当前连接处理
		if err := c.conn.ReadJSON(&input); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("websocket read:", err)
			}
			return
		}

		// 如果客户端发来的消息类型不是 "chat", 那服务端就不处理
		if input.Type != "chat" {
			c.send <- OutgoingMessage{Type: "error", Error: "unsupported message type"}
			continue
		}
		// 检查消息的基本字段是否完整
		if input.To <= 0 || input.Content == "" {
			c.send <- OutgoingMessage{Type: "error", Error: "to and content are required"}
			continue
		}
		// 检查接收者是否真的存在
		if _, ok := c.hub.store.GetUserByID(input.To); !ok {
			c.send <- OutgoingMessage{Type: "error", Error: "receiver not found"}
			continue
		}

		// 把当前用户发来的消息保存到存储里
		message, err := c.hub.store.SaveMessage(c.userID, input.To, input.Content)
		if err != nil {
			c.send <- OutgoingMessage{Type: "error", Error: "save message failed"}
			continue
		}

		// 先准备给收件人的聊天消息
		payload := OutgoingMessage{Type: "chat", Message: message}
		//告诉发件人：“发成功了”
		c.send <- OutgoingMessage{Type: "ack", Message: message}
		// 再尝试把消息实时推给收件人
		c.hub.deliver(input.To, payload)
	}
}

// 专门负责把服务端消息写到当前客户端的 WebSocket 连接里, 并定时发送 ping 保持连接存活
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case payload, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // 设置写超时
			// c.send 通道关闭
			// 这个客户端已经被注销, 不需要再发消息了
			if !ok {
				// 发一个 WebSocket 关闭消息, 然后退出函数
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			data, err := json.Marshal(payload)
			if err != nil {
				log.Println("json marshal:", err)
				continue
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// ping 无法发出 -> 结束连接
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
