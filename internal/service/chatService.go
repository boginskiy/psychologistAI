package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/ai"
	"github.com/boginskiy/psychologistAI/internal/domain/chat"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/repository"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
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
	// Пользователь не зарегистрирован. Возвращаем инициализированный пустой чат
	// Истории у незарегистрированного пользователя нет.
	if infoUser == nil {
		return &chat.Chat{}, nil
	}

	userChat, err := s.ChatRepo.Read(infoUser.UserID)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrReadChat
		}
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

	// TODO
	// Остановка тут! Ошибка! infoUser у нас есть nil, а есть только с сообщением.
	// Определиться бы не помешало.

	userChat, err := s.GetOrCreate(ctx, infoUser)
	if err != nil {
		return nil, err
	}

	// User Message
	userMessage := chat.NewMessage(userChat.ID, RoleUser, infoUser.Message)

	// AI Client
	response, err := s.AIClient.Send(ctx, infoUser.Message)
	if err != nil {
		userChat.AddMessage(*userMessage)
		s.ChatRepo.Update(userChat)
		return userChat, errs.ErrAIResponse
	}

	// AI Message
	aiMessage := chat.NewMessage(userChat.ID, RoleAssistant, response)

	userChat.AddMessage(*userMessage, *aiMessage)
	s.ChatRepo.Update(userChat)

	return userChat, nil
}
