package subscriber

import (
	"context"
	"food-delivery/component/appctx"
	restaurantstorage "food-delivery/module/restaurant/storage"
	"food-delivery/pubsub"
)

//func DecreaseLikeCountAfterUserDisLikeRestaurant(appCtx appctx.AppContext, ctx context.Context) {
//	c, _ := appCtx.GetPubSub().Subscribe(ctx, common.TopicUserLikeRestaurant)
//
//	store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
//
//	go func() {
//		defer common.AppRecover()
//		for {
//			msg := <-c
//			likeData := msg.Data().(HasRestaurantId)
//			_ = store.DecreaseLikeCount(ctx, likeData.GetRestaurantId())
//		}
//	}()
//}

func DecreaseLikeCountAfterUserDisLikeRestaurant(appCtx appctx.AppContext) consumerJob {
	return consumerJob{
		Title: "Decrease Like Count After User Disike Restaurant",
		Hdl: func(ctx context.Context, message *pubsub.Message) error {
			store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
			likeData := message.Data().(HasRestaurantId)

			return store.DecreaseLikeCount(ctx, likeData.GetRestaurantId())
		},
	}
}
