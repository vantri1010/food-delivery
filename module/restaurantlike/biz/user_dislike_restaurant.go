package rstlikebiz

import (
	"context"
	"food-delivery/common"
	"food-delivery/module/restaurantlike/model"
	"log"
)

type UserDislikeRestaurantStore interface {
	Delete(ctx context.Context, userId, restaurantId int) error
}

type DecLikedCountResStore interface {
	DecreaseLikeCount(ctx context.Context, id int) error
}

type userDislikeRestaurantBiz struct {
	store    UserDislikeRestaurantStore
	decStore DecLikedCountResStore
}

func NewUserDislikeRestaurantBiz(store UserDislikeRestaurantStore, decStore DecLikedCountResStore) *userDislikeRestaurantBiz {
	return &userDislikeRestaurantBiz{
		store: store, decStore: decStore,
	}
}

func (biz *userDislikeRestaurantBiz) DislikeRestaurant(
	ctx context.Context,
	UserId,
	restaurantId int,
) error {
	err := biz.store.Delete(ctx, UserId, restaurantId)

	if err != nil {
		return restaurantlikemodel.ErrCannotUnLikeRestaurant(err)
	}

	go func() { // goroutine for avoiding crashes
		defer common.AppRecover()

		if err := biz.decStore.DecreaseLikeCount(ctx, restaurantId); err != nil {
			// should not do this: return restaurantlikemodel.ErrCannotLikeRestaurant(err)
			log.Println(err)
		}
	}()

	return nil
}
