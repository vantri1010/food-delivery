package subscriber

import (
	"context"
	"food-delivery/component/appctx"
	restaurantstorage "food-delivery/module/restaurant/storage"
	"food-delivery/pubsub"
	"log"
)

// HasRestaurantId created to avoid restaurantlikemodel dependent
type HasRestaurantId interface {
	GetRestaurantId() int
	//GetUserId()
}

//func IncreaseLikeCountAfterUserLikeRestaurant(appCtx appctx.AppContext, ctx context.Context) {
//	c, _ := appCtx.GetPubSub().Subscribe(ctx, common.TopicUserLikeRestaurant)
//
//	store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
//
//	go func() {
//		defer common.AppRecover()
//		for {
//			msg := <-c
//			likeData := msg.Data().(HasRestaurantId)
//			_ = store.IncreaseLikeCount(ctx, likeData.GetRestaurantId())
//		}
//	}()
//}

func IncreaseLikeCountAfterUserLikeRestaurant(appCtx appctx.AppContext) consumerJob {
	return consumerJob{
		Title: "Increase Like Count After User Like Restaurant",
		Hdl: func(ctx context.Context, message *pubsub.Message) error {
			store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
			likeData := message.Data().(HasRestaurantId)

			return store.IncreaseLikeCount(ctx, likeData.GetRestaurantId())
		},
	}
}

func PushNotiUserLikeRestaurant(appCtx appctx.AppContext) consumerJob {
	return consumerJob{
		Title: "Push notification when user likes restaurant",
		Hdl: func(ctx context.Context, message *pubsub.Message) error {
			likeData := message.Data().(HasRestaurantId)
			log.Println("Push notification when user likes restaurant", likeData)

			return nil
		},
	}
}
