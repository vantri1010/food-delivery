package restaurantlikestore

import (
	"context"
	"food-delivery/common"
	"food-delivery/module/restaurantlike/model"
)

func (s *sqlStore) Create(context context.Context, data *restaurantlikemodel.Like) error {
	if err := s.db.Create(&data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
