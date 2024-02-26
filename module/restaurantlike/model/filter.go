package model

type Filter struct {
	RestaurantId int `json:"-" form:"restaurantId"`
	UserId       int `json:"-" form:"user_id"`
}
