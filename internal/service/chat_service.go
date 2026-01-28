package service

import (
	"chats/internal/domain"
	"chats/internal/repositories"
	"context"
	"strings"
)

type chatService struct {
	chatRepo repositories.ChatRepository
}

func NewChatService(chatRepo repositories.ChatRepository) ChatService {
	return &chatService{chatRepo: chatRepo}
}

func (c chatService) CreateChat(ctx context.Context, title string) (*domain.Chat, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return nil, domain.ErrInvalidInput
	}

	if len(title) > 0 {
		return nil, domain.ErrInvalidInput
	}

	chat := &domain.Chat{
		Title: title,
	}
	if err := c.chatRepo.Create(ctx, chat); err != nil {
		return nil, err
	}
	return chat, nil
}

func (c chatService) GetChat(ctx context.Context, id uint, limit int) (*domain.Chat, error) {
	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	return c.chatRepo.GetByID(ctx, id, true, limit)
}

func (c chatService) DeleteChat(ctx context.Context, id uint) error {
	return c.chatRepo.Delete(ctx, id)
}

func (c chatService) ValidateChatExists(ctx context.Context, id uint) error {
	//TODO implement me
	panic("implement me")
}
