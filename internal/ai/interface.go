package ai

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/domain/chat"
)

type AIClient interface {
	SendWithHistory(context.Context, []chat.Message, string) (string, error)
	Send(ctx context.Context, req string) (res string, err error)
}
