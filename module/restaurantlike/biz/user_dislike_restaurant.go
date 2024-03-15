package rstlikebiz

import (
	"context"
	"food-delivery/common"
	"food-delivery/module/restaurantlike/model"
	"food-delivery/pubsub"
	"log"
)

type UserDislikeRestaurantStore interface {
	Delete(ctx context.Context, userId, restaurantId int) error
}

//type DecLikedCountResStore interface {
//	DecreaseLikeCount(ctx context.Context, id int) error
//}

type userDislikeRestaurantBiz struct {
	store UserDislikeRestaurantStore
	//decStore DecLikedCountResStore
	ps pubsub.Pubsub
}

func NewUserDislikeRestaurantBiz(
	store UserDislikeRestaurantStore,
	//decStore DecLikedCountResStore,
	ps pubsub.Pubsub,
) *userDislikeRestaurantBiz {
	return &userDislikeRestaurantBiz{
		store: store,
		//decStore: decStore,
		ps: ps,
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

	if err := biz.ps.Publish(
		ctx,
		common.TopicUserDisLikeRestaurant,
		pubsub.NewMessage(&restaurantlikemodel.Like{RestaurantId: restaurantId}),
	); err != nil {
		log.Println(err)
	}

	// Side effective goroutine for avoiding crashes
	//j := asyncjob.NewJob(func(ctx context.Context) error {
	//	return biz.decStore.DecreaseLikeCount(ctx, restaurantId)
	//})
	//
	//if err := asyncjob.NewGroup(true, j).Run(ctx); err != nil {
	//	log.Println(err)
	//}

	return nil
}
