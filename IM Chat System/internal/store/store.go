package store

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	Nickname     string    `json:"nickname"`
	CreatedAt    time.Time `json:"created_at"`
}

type Message struct {
	ID        int64     `json:"id"`
	FromID    int64     `json:"from_id"`
	ToID      int64     `json:"to_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	mu            sync.RWMutex
	path          string
	NextUserID    int64     `json:"next_user_id"`
	NextMessageID int64     `json:"next_message_id"`
	Users         []User    `json:"users"`
	Messages      []Message `json:"messages"`
}

// 根据给定的文件路径, 打开并初始化一个 Store
// 如果文件里已经有历史数据, 就把数据读出来恢复到内存里
func Open(path string) (*Store, error) {
	s := &Store{
		path:          path,
		NextUserID:    1,
		NextMessageID: 1,
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return s, nil
	}
	// 把 data 中的 JSON 数据解析到 s 中
	// Unmarshal 把其他格式的数据解析成 Go 的变量
	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}

	s.path = path
	if s.NextUserID == 0 {
		s.NextUserID = int64(len(s.Users)) + 1
	}
	if s.NextMessageID == 0 {
		s.NextMessageID = int64(len(s.Messages)) + 1
	}

	return s, nil
}

func (s *Store) CreateUser(username, passwordHash, nickname string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.findUserByUsernameLocked(username) != nil {
		return User{}, errors.New("username already exists")
	}
	if nickname == "" {
		nickname = username
	}

	user := User{
		ID:           s.NextUserID,
		Username:     username,
		PasswordHash: passwordHash,
		Nickname:     nickname,
		CreatedAt:    time.Now(),
	}
	s.NextUserID++
	s.Users = append(s.Users, user)
	return user, s.saveLocked()
}

// 按用户名查用户, 找到返回用户和 true, 找不到返回空用户和 false
func (s *Store) GetUserByUsername(username string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user := s.findUserByUsernameLocked(username)
	if user == nil {
		return User{}, false
	}
	return *user, true
}

// 按用户 ID 查人, 找到返回用户和 true, 找不到返回空用户和 false
func (s *Store) GetUserByID(id int64) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.Users {
		if user.ID == id {
			return user, true
		}
	}
	return User{}, false
}

// 读取所有用户, 但返回前先把密码哈希去掉, 避免泄露敏感信息
func (s *Store) ListUsers() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]User, len(s.Users))
	copy(users, s.Users)
	for i := range users {
		users[i].PasswordHash = ""
	}
	return users
}

// 把一条新聊天消息保存到 Store 里, 并写入 data.json, 最后把这条消息返回出去
func (s *Store) SaveMessage(fromID, toID int64, content string) (Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	message := Message{
		ID:        s.NextMessageID,
		FromID:    fromID,
		ToID:      toID,
		Content:   content,
		CreatedAt: time.Now(),
	}
	s.NextMessageID++
	s.Messages = append(s.Messages, message)
	return message, s.saveLocked()
}

// 查询两个人的历史聊天消息
func (s *Store) Conversation(userID, peerID, afterID int64) []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var messages []Message
	for _, message := range s.Messages {
		isPair := (message.FromID == userID && message.ToID == peerID) ||
			(message.FromID == peerID && message.ToID == userID)
		if isPair && message.ID > afterID {
			messages = append(messages, message)
		}
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].ID < messages[j].ID
	})
	return messages
}

// 从消息库里找出某个用户还没拉取过的新消息, 并按时间顺序返回
func (s *Store) OfflineMessages(userID, afterID int64) []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var messages []Message
	for _, message := range s.Messages {
		if message.ToID == userID && message.ID > afterID {
			messages = append(messages, message)
		}
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].ID < messages[j].ID
	})
	return messages
}

// 根据用户名在用户列表里查找用户
func (s *Store) findUserByUsernameLocked(username string) *User {
	for i := range s.Users {
		if s.Users[i].Username == username {
			return &s.Users[i]
		}
	}
	return nil
}

// 把当前 Store 里的数据保存到 data.json 文件
func (s *Store) saveLocked() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}
