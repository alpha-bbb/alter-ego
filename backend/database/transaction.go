package database

import (
	"gorm.io/gorm"
)

type ITransaction interface {
	StartTransaction(function func(tx any) error) error
}

func NewGormTransaction(db *gorm.DB) ITransaction {
	return &GormTransaction{
		db: db,
	}
}

type GormTransaction struct {
	db *gorm.DB
}

func (t *GormTransaction) StartTransaction(function func(tx any) error) error {
	return t.db.Transaction(func(gormTx *gorm.DB) error {
		return function(gormTx)
	})
}
