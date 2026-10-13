package service

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/ai"
	"github.com/boginskiy/psychologistAI/internal/domain/chat"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/repository"
	"github.com/google/uuid"
)

const (
	AIConsumer = "user"
	AIProducer = "assistant"
)

type ChatServ struct {
	AIClient ai.AIClient
	ChatRepo repository.ChatRepo
}

func NewChatServ(ctx context.Context, aiClient ai.AIClient, chatRepo repository.ChatRepo) *ChatServ {
	return &ChatServ{
		AIClient: aiClient,
		ChatRepo: chatRepo,
	}
}

func (s *ChatServ) GetOrCreate(ctx context.Context, infoUser *dto.InfoUser) (*chat.Chat, error) {
	// Это костылек мотылек
	if infoUser == nil {
		return &chat.Chat{}, nil
	}

	// GuestChat
	if infoUser.UserID == uuid.Nil {
		guestChat := chat.NewChat(infoUser)
		return guestChat, nil
	}

	// UserChat
	userChat, err := s.ChatRepo.Read(infoUser.UserID)
	if err != nil {
		// if !errors.Is(err, sql.ErrNoRows) {
		// 	return nil, errs.ErrReadChat
		// }
		userChat = chat.NewChat(infoUser)
	}
	return userChat, nil
}

func (s *ChatServ) Send(ctx context.Context, infoUser *dto.InfoUser) (*chat.Chat, error) {
	// Если это не авторизованный пользователь, то у него есть 5 запросов
	// Данные по таким пользователям будем хранить в отдельно таблице, или Redis
	if infoUser.SessionID == "" {
		cnt := 0
		cnt++
	}

	// User Chat
	userChat, err := s.GetOrCreate(ctx, infoUser)
	if err != nil {
		return nil, err
	}

	// User Message
	userMessage := chat.NewMessage(userChat.ID, AIConsumer, infoUser.Message)
	userChat.AddMessage(userMessage)

	// AI Client
	response, err := s.AIClient.Send(ctx, infoUser.Message)
	if err != nil {
		s.ChatRepo.Update(userChat)
		return userChat, errs.ErrAIResponse
	}

	// AI Message
	aiMessage := chat.NewMessage(userChat.ID, AIProducer, response)

	userChat.AddMessage(aiMessage)
	s.ChatRepo.Update(userChat)

	return userChat, nil
}
