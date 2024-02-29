package rstlikebiz

import (
	"context"
	"food-delivery/common"
	"food-delivery/module/restaurantlike/model"
	"log"
)

type UserLikeRestaurantStore interface {
	Create(ctx context.Context, data *restaurantlikemodel.Like) error
}

type IncLikedCountResStore interface {
	IncreaseLikeCount(ctx context.Context, id int) error
}

type userLikeRestaurantBiz struct {
	store    UserLikeRestaurantStore
	IncStore IncLikedCountResStore
}

func NewUserLikeRestaurantBiz(store UserLikeRestaurantStore, IncStore IncLikedCountResStore) *userLikeRestaurantBiz {
	return &userLikeRestaurantBiz{store: store, IncStore: IncStore}
}

func (biz *userLikeRestaurantBiz) LikeRestaurant(ctx context.Context, data *restaurantlikemodel.Like) error {
	err := biz.store.Create(ctx, data)

	if err != nil {
		return restaurantlikemodel.ErrCannotLikeRestaurant(err)
	}

	go func() { // goroutine for avoiding crashes
		defer common.AppRecover()

		if err := biz.IncStore.IncreaseLikeCount(ctx, data.RestaurantId); err != nil {
			// should not do this: return restaurantlikemodel.ErrCannotLikeRestaurant(err)
			log.Println(err)
		}
	}()

	return nil
}
