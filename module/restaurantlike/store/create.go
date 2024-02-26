package store

import (
	"context"
	"food-delivery/common"
	"food-delivery/module/restaurantlike/model"
)

func (s *sqlStore) Create(ctx context.Context, data *model.Like) error {
	db := s.db

	if err := db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
