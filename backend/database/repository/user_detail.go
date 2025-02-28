package repository

import (
	"errors"

	"github.com/alpha-bbb/alter-ego/backend/database/model"
	"github.com/alpha-bbb/alter-ego/backend/entity"
	"gorm.io/gorm"
)

type IUserDetailRepository interface {
	CreateInTx(tx any, entity entity.UserDetail) error
	FindByUserID(userID string) (entity.UserDetail, error)
}

type UserDetailRepository struct {
	db *gorm.DB
}

func NewUserDetailRepository(db *gorm.DB) IUserDetailRepository {
	return &UserDetailRepository{db: db}
}

func (r *UserDetailRepository) CreateInTx(tx any, entity entity.UserDetail) error {
	txAsserted, ok := tx.(*gorm.DB)
	if !ok {
		return ErrInvalidTransaction
	}

	userDetail := &model.UserDetail{
		UserDetailID: entity.UserDetailID,
		UserID:       entity.User.UserID,
	}
	return txAsserted.
		Create(userDetail).
		Error
}

func (r *UserDetailRepository) FindByUserID(userID string) (entity.UserDetail, error) {
	var userDetail model.UserDetail
	err := r.db.Model(&model.UserDetail{}).
		Where("user_id = ?", userID).
		Preload("User").
		First(&userDetail).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.UserDetail{}, ErrNotFound
		}
		return entity.UserDetail{}, err
	}
	return userDetail.Entity(), nil
}
