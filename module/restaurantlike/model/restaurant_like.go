package restaurantlikemodel

import (
	"fmt"
	"food-delivery/common"
	"time"
)

const EntityName = "UserLikeRestaurant"

type Like struct {
	RestaurantId int `json:"restaurant_id" gorm:"column:restaurant_id"`
	UserId       int `json:"user_id" gorm:"column:user_id"`
	//CreatedAt    *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	CreatedAt *time.Time         `json:"created_at,omitempty" gorm:"column:created_at;type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP"`
	User      *common.SimpleUser `json:"user,omitempty" gorm:"<-:false"` // use to preload users. read only for avoid insert into associated table user when insert into Like
}

func (Like) TableName() string {
	return "restaurant_likes"
}

func (l *Like) GetRestaurantId() int {
	return l.RestaurantId
}

func ErrCannotLikeRestaurant(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("cannot like this restanrant"),
		fmt.Sprintf("ErrCannotLikeRestaurant"))
}

func ErrCannotUnLikeRestaurant(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("cannot dislike this restanrant"),
		fmt.Sprintf("ErrCannotDislikeRestaurant"))
}
