package restaurantlikestore

import (
	"context"
	"food-delivery/common"
	"food-delivery/module/restaurantlike/model"
	"github.com/btcsuite/btcd/btcutil/base58"
	"time"
)

func (s *sqlStore) GetRestaurantLikes(ctx context.Context, ids []int) (map[int]int, error) {
	result := make(map[int]int)

	type sqlData struct {
		RestaurantId int `gorm:"column:restaurant_id"`
		LikeCount    int `gorm:"column:count"`
	}

	var listLike []sqlData

	if err := s.db.Table(restaurantlikemodel.Like{}.TableName()).
		Select("restaurant_id, count(restaurant_id) as count").
		Where("restaurant_id in (?)", ids).
		Group("restaurant_id").
		Find(&listLike).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for _, item := range listLike {
		result[item.RestaurantId] = item.LikeCount
	}

	return result, nil
}

func (s *sqlStore) GetUsersLikeRestaurant(
	ctx context.Context,
	conditions map[string]interface{},
	filter *restaurantlikemodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]common.SimpleUser, error) {
	var likes []restaurantlikemodel.Like

	db := s.db

	db = db.Table(restaurantlikemodel.Like{}.TableName()).Where(conditions)

	if filter.RestaurantId > 0 {
		db.Where("restaurant_id = ?", filter.RestaurantId)
	}

	if err := db.Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	// there are 2 type of paging. by cursor and by offset depend on the user's input
	// cursor is encoded restaurantId copy from list restaurants api
	if v := paging.FakeCursor; v != "" {
		//in this case. cursor is the field timestamp when user like restaurant. due to we have multiple primary keys
		timeCreated, err := time.Parse(time.DateTime, string(base58.Decode(v)))

		if err != nil {
			return nil, common.ErrDB(err)
		}

		db = db.Where("created_at <= ?", timeCreated.Format(time.DateTime))
	} else {
		// offsets the query by the number of records to skip (based on the current page number)
		// and sets the limit to the number of records to fetch (based on the paging.Limit field).
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.Preload("User").
		Limit(paging.Limit + 1).
		Order("created_at desc").
		Find(&likes).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	// handle the case where we in the last page
	if len(likes) <= paging.Limit {
		paging.NextCursor = ""
	} else {
		paging.NextCursor = base58.Encode([]byte(likes[len(likes)-1].CreatedAt.Format(time.DateTime)))
		likes = likes[:len(likes)-1]
	}

	users := make([]common.SimpleUser, len(likes))

	for i, item := range likes {
		users[i] = *likes[i].User
		users[i].CreatedAt = item.CreatedAt
		users[i].UpdatedAt = nil
	}

	return users, nil
}
