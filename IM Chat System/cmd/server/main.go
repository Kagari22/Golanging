package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"IM_Chat_System/internal/auth"
	"IM_Chat_System/internal/chat"
	"IM_Chat_System/internal/store"
)

/*
启动方式 (powershell)
$env:IM_ADDR=":8081"
go run ./cmd/server
*/

type app struct {
	store  *store.Store
	hub    *chat.Hub
	secret string
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

func main() {
	secret := getenv("IM_JWT_SECRET", "change-me-in-real-project")
	dataPath := getenv("IM_DATA_PATH", "data.json")

	// 打开 data.json，准备保存用户和消息
	s, err := store.Open(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	// 创建 app，保存全局依赖：store 和 secret
	a := &app{
		store:  s,
		secret: secret,
	}
	// 创建聊天 Hub，用来管理 WebSocket 在线用户和消息转发
	a.hub = chat.NewHub(s, secret)

	// 注册路由并启动 HTTP 服务器
	// mux 是路由器
	mux := http.NewServeMux() // 根据请求路径和请求方法，把请求分发给对应处理函数
	// 给不同的路径注册不同的处理函数
	mux.HandleFunc("POST /api/register", a.handleRegister) // 如果请求是 POST /api/register 就执行 a.handleRegister
	mux.HandleFunc("POST /api/login", a.handleLogin)
	// 带了 a.withAuth(...) 表示这些接口需要登录后才能访问
	mux.HandleFunc("GET /api/me", a.withAuth(a.handleMe)) // 获取当前登录用户信息
	mux.HandleFunc("GET /api/users", a.withAuth(a.handleUsers)) // 获取其他用户列表
	mux.HandleFunc("GET /api/messages", a.withAuth(a.handleMessages)) // 获取和某个人的历史消息
	mux.HandleFunc("GET /api/offline", a.withAuth(a.handleOfflineMessages)) // 获取离线消息
	mux.HandleFunc("GET /ws", a.hub.ServeWS) // WebSocket 接口
	mux.Handle("/", http.FileServer(http.Dir("web"))) // 静态文件服务

	addr := getenv("IM_ADDR", ":8080") 
	log.Println("server listening on", addr)
	// 在 addr 这个地址监听请求
	// 收到请求后交给 logRequest(mux) 处理
	// 如果启动失败，就打印错误并退出
	log.Fatal(http.ListenAndServe(addr, logRequest(mux)))
}

// 注册接口的处理函数
func (a *app) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	// 从 HTTP 请求体 r.Body 里读取 JSON
	// 然后解析到 req 这个变量里
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Nickname = strings.TrimSpace(req.Nickname)
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password are required")
		return
	}
	if len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "hash password failed")
		return
	}

	user, err := a.store.CreateUser(req.Username, passwordHash, req.Nickname)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	// 注册成功, 返回 201 状态码, 并把新用户的安全信息以 JSON 形式返回给客户端
	writeJSON(w, http.StatusCreated, map[string]any{
		"user": toUserResponse(user),
	})
}

func (a *app) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	// 用用户名去 store 里找用户
	user, ok := a.store.GetUserByUsername(strings.TrimSpace(req.Username))
	if !ok || !auth.CheckPassword(req.Password, user.PasswordHash) {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	token, err := auth.GenerateToken(a.secret, user.ID, user.Username, 7*24*time.Hour)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "generate token failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  toUserResponse(user),
	})
}

// 获取当前登录用户信息
func (a *app) handleMe(w http.ResponseWriter, r *http.Request, claims auth.Claims) {
	user, ok := a.store.GetUserByID(claims.UserID)
	if !ok {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": toUserResponse(user)})
}

// 获取用户列表
func (a *app) handleUsers(w http.ResponseWriter, r *http.Request, claims auth.Claims) {
	users := a.store.ListUsers() // 从 store 里拿到所有用户
	resp := make([]userResponse, 0, len(users))
	for _, user := range users {
		if user.ID != claims.UserID {
			resp = append(resp, toUserResponse(user))
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"users": resp})
}

// 获取和某个用户的聊天记录
func (a *app) handleMessages(w http.ResponseWriter, r *http.Request, claims auth.Claims) {
	peerID, err := strconv.ParseInt(r.URL.Query().Get("peer_id"), 10, 64)
	if err != nil || peerID <= 0 {
		writeError(w, http.StatusBadRequest, "peer_id is required")
		return
	}
	afterID, _ := strconv.ParseInt(r.URL.Query().Get("after_id"), 10, 64)
	// 查询当前登录用户 claims.UserID 和 peerID 这个用户之间, ID 大于 afterID 的聊天消息
	messages := a.store.Conversation(claims.UserID, peerID, afterID)
	writeJSON(w, http.StatusOK, map[string]any{"messages": messages})
}

// 获取离线消息
func (a *app) handleOfflineMessages(w http.ResponseWriter, r *http.Request, claims auth.Claims) {
	afterID, _ := strconv.ParseInt(r.URL.Query().Get("after_id"), 10, 64)
	messages := a.store.OfflineMessages(claims.UserID, afterID)
	writeJSON(w, http.StatusOK, map[string]any{"messages": messages})
}

// 登录鉴权中间件
// 在真正执行接口之前, 先检查用户有没有带 token, token 是否有效
func (a *app) withAuth(next func(http.ResponseWriter, *http.Request, auth.Claims)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.BearerToken(r.Header.Get("Authorization"))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "missing token")
			return
		}

		claims, err := auth.ParseToken(a.secret, token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		next(w, r, claims)
	}
}

// 把内部用户对象转换成返回给前端的用户对象
func toUserResponse(user store.User) userResponse {
	return userResponse{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	// 告诉客户端：我返回的是 JSON 数据, 编码是 utf-8
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status) // 设置 HTTP 状态码
	json.NewEncoder(w).Encode(data) // 把数据转成 JSON 返回
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// 安全地读取环境变量
func getenv(key, fallback string) string {
	value := os.Getenv(key) // 尝试读取环境变量 key 的值
	if value == "" {
		return fallback
	}
	return value
}

// 每次有 HTTP 请求进来时, 先打印请求方法和路径, 然后再交给真正的路由处理
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
