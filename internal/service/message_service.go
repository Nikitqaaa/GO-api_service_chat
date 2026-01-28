package service

import (
	"chats/internal/domain"
	"chats/internal/repositories"
	"context"
	"strings"
)

type messagesService struct {
	messagesRepo repositories.MessageRepository
	chatService  ChatService
}

func NewMessagesService(messagesRepo repositories.MessageRepository, chatService ChatService) MessagesService {
	return &messagesService{
		messagesRepo: messagesRepo,
		chatService:  chatService,
	}
}

func (m messagesService) CreateMessages(ctx context.Context, chatID uint, text string) (*domain.Message, error) {
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

	if err := m.messagesRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}
