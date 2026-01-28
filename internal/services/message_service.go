package services

import (
	"chats/internal/domain"
	"chats/internal/repositories"
	"context"
	"strings"
)

type messageService struct {
	messageRepo repositories.MessageRepository
	chatService ChatService
}

func NewMessageService(messageRepo repositories.MessageRepository, chatService ChatService) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		chatService: chatService,
	}
}

func (m messageService) CreateMessage(ctx context.Context, chatID uint, text string) (*domain.Message, error) {
	if err := m.chatService.ValidateChatExists(ctx, chatID); err != nil {
		return nil, err
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return nil, domain.ErrInvalidInput
	}

	message := &domain.Message{
		ChatID: chatID,
		Text:   text,
	}

	if err := m.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}
