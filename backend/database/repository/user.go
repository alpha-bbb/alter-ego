package repository

import (
	"errors"
	"fmt"

	"github.com/alpha-bbb/alter-ego/backend/database/model"
	"github.com/alpha-bbb/alter-ego/backend/entity"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateInTx(tx any, entity entity.User) error
	FindByID(userID string) (entity.User, error)
	FindByAccount(platform, accountID string) (entity.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateInTx(tx any, entity entity.User) error {
	txAsserted, ok := tx.(*gorm.DB)
	if !ok {
		return ErrInvalidTransaction
	}

	model := &model.User{
		UserID: entity.UserID,
	}
	return txAsserted.Create(model).Error
}

func (r *UserRepository) FindByID(userID string) (entity.User, error) {
	var user model.User
	err := r.db.Model(&model.User{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, ErrNotFound
		}
		return entity.User{}, fmt.Errorf("failed to retrieve user: %w", err)
	}
	return user.Entity(), nil
}

func (r *UserRepository) FindByAccount(platform, accountID string) (entity.User, error) {
	// XXX: 現在はLINEしか対応していないためplatformはLINE決め打ち
	// 複数ある場合はswitch文で分岐
	var user model.User
	err := r.db.Model(&model.User{}).
		Joins("JOIN user_account_lines ON user_account_lines.user_id = users.user_id").
		Where("user_account_lines.line_id = ?", accountID).
		Order("user_account_lines.created_at DESC").
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, ErrNotFound
		}
		return entity.User{}, fmt.Errorf("failed to retrieve user: %w", err)
	}
	return user.Entity(), nil
}
