package userstore

import (
	"context"
	"food-delivery/common"
	usermodel "food-delivery/module/user/model"
)

func (s *sqlStore) CreateUser(ctx context.Context, data *usermodel.UserCreate) error {
	db := s.db.Begin() //open a connection

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		db.Rollback() //avoid too many connections to DB
		return common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}
