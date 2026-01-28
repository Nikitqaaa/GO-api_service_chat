package repositories

import (
	"chats/internal/domain"
	"context"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *domain.Chat) error
	GetByID(ctx context.Context, id uint, withMessage bool, limit int) (*domain.Chat, error)
	Delete(ctx context.Context, id uint) error
	Exists(ctx context.Context, id uint) (bool, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	GetByChatID(ctx context.Context, chatID uint, limit int) ([]domain.Message, error)
}
