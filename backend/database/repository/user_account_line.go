package repository

import (
	"errors"

	"github.com/alpha-bbb/alter-ego/backend/database/model"
	"github.com/alpha-bbb/alter-ego/backend/entity"
	"gorm.io/gorm"
)

type IUserAccountLineRepository interface {
	CreateInTx(tx any, entity entity.UserAccountLine) error
	FindByUserID(userID string) (entity.UserAccountLine, error)
}

type UserAccountLineRepository struct {
	db *gorm.DB
}

func NewUserAccountLineRepository(db *gorm.DB) IUserAccountLineRepository {
	return &UserAccountLineRepository{db: db}
}

func (r *UserAccountLineRepository) CreateInTx(tx any, entity entity.UserAccountLine) error {
	txAsserted, ok := tx.(*gorm.DB)
	if !ok {
		return ErrInvalidTransaction
	}

	model := &model.UserAccountLine{
		UserAccountLineID: entity.UserAccountLineID,
		UserID:            entity.User.UserID,
		LineID:            entity.LineID,
	}
	return txAsserted.
		Create(model).
		Error
}

func (r *UserAccountLineRepository) FindByUserID(userID string) (entity.UserAccountLine, error) {
	var userAccountLine model.UserAccountLine
	err := r.db.Model(&model.UserAccountLine{}).
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&userAccountLine).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.UserAccountLine{}, ErrNotFound
		}
		return entity.UserAccountLine{}, err
	}
	return userAccountLine.Entity(), nil
}
