package restaurantbiz

import (
	"context"
	"food-delivery/common"
	"food-delivery/module/restaurant/model"
	"log"
)

type ListRestaurantStore interface {
	ListDataWithCondition(
		context context.Context,
		filter *restaurantmodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]restaurantmodel.Restaurant, error)
}

type LikeRestaurantStore interface {
	GetRestaurantLikes(ctx context.Context, ids []int) (map[int]int, error)
}

type listRestaurantBiz struct {
	store     ListRestaurantStore
	likeStore LikeRestaurantStore
}

func NewListRestaurantBiz(store ListRestaurantStore, likeStore LikeRestaurantStore) *listRestaurantBiz {
	return &listRestaurantBiz{store: store, likeStore: likeStore}
}

func (biz *listRestaurantBiz) ListRestaurant(
	context context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging,
) ([]restaurantmodel.Restaurant, error) {
	result, err := biz.store.ListDataWithCondition(context, filter, paging, "User")

	if err != nil {
		return nil, err
	}

	ids := make([]int, len(result))

	for i := range ids {
		ids[i] = result[i].Id
	}

	likeMap, err := biz.likeStore.GetRestaurantLikes(context, ids)

	if err != nil {
		log.Println(err)
		return result, nil
	}

	for i, item := range result {
		result[i].LikedCount = likeMap[item.Id]
	}

	return result, nil
}
