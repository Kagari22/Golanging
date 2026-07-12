package memory

import "go-agent/llm"

type Store struct {
	limit    int
	messages []llm.Message
}

func NewStore(limit int) *Store {
	return &Store{limit: limit}
}

// 往记忆里追加一条消息, 并在超过上限时, 只保留最近的几条
func (s *Store) Add(message llm.Message) {
	s.messages = append(s.messages, message)
	if s.limit > 0 && len(s.messages) > s.limit {
		s.messages = append([]llm.Message(nil), s.messages[len(s.messages) - s.limit:]...)
	}
}

// 把当前记忆里的消息返回出去, 但返回的是一份副本, 不是原始切片本体
func (s *Store) Messages() []llm.Message {
	return append([]llm.Message(nil), s.messages...)
}
