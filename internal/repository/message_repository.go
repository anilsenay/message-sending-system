package repository

import (
	"context"

	"github.com/anilsenay/message-sending-system/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type db interface {
	GetConnection() *gorm.DB
	Close()
}

type MessageRepository struct {
	db db
}

func NewMessageRepository(db db) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) RetrieveAll(ctx context.Context, filters model.DbFilters[model.Message]) ([]model.Message, error) {
	var data []model.Message
	result := r.db.GetConnection().WithContext(ctx).Model(&model.Message{})

	if filters.Model != nil {
		result = result.Where(*filters.Model)
	}
	if filters.Limit != nil {
		result = result.Limit(*filters.Limit)
	}
	if filters.Offset != nil {
		result = result.Offset(*filters.Offset)
	}
	if filters.Order != nil {
		result = result.Order(*filters.Order)
	}

	result = result.Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func (r *MessageRepository) Update(ctx context.Context, model *model.Message, updates map[string]interface{}) error {
	result := r.db.GetConnection().WithContext(ctx).Clauses(clause.Returning{}).Model(&model).Where(model).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MessageRepository) Create(ctx context.Context, model *model.Message) error {
	result := r.db.GetConnection().WithContext(ctx).Model(&model).Create(model)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MessageRepository) RetrieveMessagesForProcess(ctx context.Context, limit int) ([]model.Message, error) {
	var data []model.Message
	err := r.db.GetConnection().Transaction(func(tx *gorm.DB) error {
		result := tx.
			WithContext(ctx).
			Model(&model.Message{}).
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where(model.Message{Status: model.MESSAGE_UNSENT}).
			Limit(limit).
			Order("created_at ASC").
			Find(&data)
		if result.Error != nil {
			return result.Error
		}

		if len(data) == 0 {
			return nil
		}

		ids := []int{}
		for _, message := range data {
			ids = append(ids, message.Id)
		}

		result = tx.Model(&data).Clauses(clause.Returning{}).Where("id IN ?", ids).Update("status", model.MESSAGE_PROCESSING)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return data, err
}
