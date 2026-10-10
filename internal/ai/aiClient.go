package ai

import (
	"context"
	"fmt"

	"github.com/boginskiy/psychologistAI/internal/domain/chat"
)

var Mess = map[string]string{
	"0": "Мои дела отлично, а как твои?",
	"1": "Мой день все в битах и байтах, ничего не меняется!",
	"2": "Я не человек, мне никакого удовольствия в этом",
	"3": "Я с удовольствием, только генератор купить надо",
}

type AIManager struct {
}

func NewAIManager(ctx context.Context) *AIManager {
	return &AIManager{}
}

func (c *AIManager) Send(ctx context.Context, req string) (res string, err error) {
	if msg, ok := Mess[req]; ok {
		return msg, nil
	}
	return "", fmt.Errorf("bad message")
}

func (c *AIManager) SendWithHistory(ctx context.Context, history []chat.Message, msg string) (res string, err error) {
	return "", nil
}
