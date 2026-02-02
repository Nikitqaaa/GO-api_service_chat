package repositories

import (
	"chats/internal/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{
		db: db,
	}
}

func (c chatRepository) Create(ctx context.Context, chat *domain.Chat) error {
	return c.db.WithContext(ctx).Create(chat).Error
}

func (c chatRepository) GetByID(ctx context.Context, id uint, withMessage bool, limit int) (*domain.Chat, error) {
	var chat domain.Chat

	query := c.db.WithContext(ctx).Model(&domain.Chat{}).Where("id = ?", id)

	if withMessage {
		query = query.Preload("Message", func(db *gorm.DB) *gorm.DB {
			return db.Order("messages.id DESC").Limit(limit)
		})
	}

	if err := query.First(&chat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &chat, nil
}

func (c chatRepository) Delete(ctx context.Context, id uint) error {
	result := c.db.WithContext(ctx).Delete(&domain.Chat{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (c chatRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64

	err := c.db.WithContext(ctx).Model(&domain.Chat{}).
		Where("id = ?", id).
		Count(&count).
		Error

	return count > 0, err
}
