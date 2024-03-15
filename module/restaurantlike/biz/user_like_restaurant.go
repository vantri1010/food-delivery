package rstlikebiz

import (
	"context"
	"food-delivery/common"
	"food-delivery/module/restaurantlike/model"
	"food-delivery/pubsub"
	"log"
)

type UserLikeRestaurantStore interface {
	Create(ctx context.Context, data *restaurantlikemodel.Like) error
}

//type IncLikedCountResStore interface {
//	IncreaseLikeCount(ctx context.Context, id int) error
//}

type userLikeRestaurantBiz struct {
	store UserLikeRestaurantStore
	//IncStore IncLikedCountResStore
	ps pubsub.Pubsub
}

func NewUserLikeRestaurantBiz(
	store UserLikeRestaurantStore,
	//IncStore IncLikedCountResStore,
	ps pubsub.Pubsub,
) *userLikeRestaurantBiz {
	return &userLikeRestaurantBiz{
		store: store,
		//IncStore: IncStore,
		ps: ps,
	}
}

func (biz *userLikeRestaurantBiz) LikeRestaurant(ctx context.Context, data *restaurantlikemodel.Like) error {
	err := biz.store.Create(ctx, data)

	if err != nil {
		return restaurantlikemodel.ErrCannotLikeRestaurant(err)
	}

	// Send message
	if err := biz.ps.Publish(ctx, common.TopicUserLikeRestaurant, pubsub.NewMessage(data)); err != nil {
		log.Println(err)
	}

	// Side effective goroutine for avoiding crashes
	//j := asyncjob.NewJob(func(ctx context.Context) error {
	//	return biz.IncStore.IncreaseLikeCount(ctx, data.RestaurantId)
	//})
	//
	//if err := asyncjob.NewGroup(true, j).Run(ctx); err != nil {
	//	log.Println(err)
	//}

	return nil
}
