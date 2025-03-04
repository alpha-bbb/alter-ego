package repository

import (
	"errors"
	"fmt"

	"github.com/alpha-bbb/alter-ego/backend/database/model"
	"github.com/alpha-bbb/alter-ego/backend/entity"
	"gorm.io/gorm"
)

type IStateCountRepository interface {
	CreateInTx(tx any, entity entity.StateCount) error
	GetAll() ([]entity.StateCount, error)
	FindByUserID(userID string) (entity.StateCount, error)
	IncrementCountInTx(tx any, userID string) error
	Clear() error
}

type StateCountRepository struct {
	db *gorm.DB
}

func NewStateCountRepository(db *gorm.DB) IStateCountRepository {
	return &StateCountRepository{db: db}
}

func (r *StateCountRepository) CreateInTx(tx any, entity entity.StateCount) error {
	txAsserted, ok := tx.(*gorm.DB)
	if !ok {
		return ErrInvalidTransaction
	}

	stateCount := &model.StateCount{
		StateCountID: entity.StateCountID,
		UserID:       entity.User.UserID,
		Count:        entity.Count,
	}
	return txAsserted.
		Create(stateCount).
		Error
}

func (r *StateCountRepository) GetAll() ([]entity.StateCount, error) {
	var stateCounts []model.StateCount
	if err := r.db.Model(&model.StateCount{}).
		Preload("User").
		Find(&stateCounts).
		Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve state counts: %w", err)
	}

	var entities []entity.StateCount
	for _, stateCount := range stateCounts {
		entities = append(entities, stateCount.Entity())
	}
	return entities, nil
}

func (r *StateCountRepository) FindByUserID(userID string) (entity.StateCount, error) {
	var stateCount model.StateCount
	if err := r.db.Model(&model.StateCount{}).
		Preload("User").
		Where("user_id", userID).
		Order("updated_at DESC").
		First(&stateCount).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.StateCount{}, ErrNotFound
		}
		return entity.StateCount{}, fmt.Errorf("failed to retrieve state count: %w", err)
	}
	return stateCount.Entity(), nil
}

func (r *StateCountRepository) IncrementCountInTx(tx any, userID string) error {
	txAsserted, ok := tx.(*gorm.DB)
	if !ok {
		return ErrInvalidTransaction
	}

	if err := txAsserted.Model(&model.StateCount{}).
		Where("user_id", userID).
		Order("updated_at DESC").
		Update("count", gorm.Expr("count + ?", 1)).
		Error; err != nil {
		return fmt.Errorf("failed to increment count: %w", err)
	}
	return nil
}

func (r *StateCountRepository) Clear() error {
	if err := r.db.Exec("TRUNCATE TABLE state_counts").
		Error; err != nil {
		return fmt.Errorf("failed to reset state count: %w", err)
	}
	return nil
}
