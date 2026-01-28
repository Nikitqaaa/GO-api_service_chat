package service

import (
	"chats/internal/domain"
	"context"
)

type ChatService interface {
	CreateChat(ctx context.Context, title string) (*domain.Chat, error)
	GetChat(ctx context.Context, id uint, limit int) (*domain.Chat, error)
	DeleteChat(ctx context.Context, id uint) error
	ValidateChatExists(ctx context.Context, id uint) error
}

type MessagesService interface {
	CreateMessages(ctx context.Context, chatID uint, message string) (*domain.Message, error)
}
