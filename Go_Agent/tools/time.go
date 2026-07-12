package tools

import (
	"context"
	"time"
)

type TimeTool struct{}

func NewTimeTool() *TimeTool {
	return &TimeTool{}
}

func (t *TimeTool) Name() string {
	return "time"
}

func (t *TimeTool) Description() string {
	return "Returns the current local time."
}

func (t *TimeTool) Run(_ context.Context, _ string) (string, error) {
	return time.Now().Format(time.RFC3339), nil
}
