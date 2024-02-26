package store

import "gorm.io/gorm"

type sqlStore struct {
	db *gorm.DB
}

func NewSQLstore(db *gorm.DB) *sqlStore {
	return &sqlStore{db: db}
}
