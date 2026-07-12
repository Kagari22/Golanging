package tools

import "context"

type Tool interface {
	Name() string
	Description() string
	Run(ctx context.Context, input string) (string, error)
}
