package repositories

import (
	"chats/internal/domain"
	"context"

	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (m messageRepository) Create(ctx context.Context, message *domain.Message) error {
	return m.db.WithContext(ctx).Create(message).Error
}

func (m messageRepository) GetByChatID(ctx context.Context, chatID uint, limit int) ([]domain.Message, error) {
	var messages []domain.Message

	err := m.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}
