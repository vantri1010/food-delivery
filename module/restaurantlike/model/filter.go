package restaurantlikemodel

// when we need to know how many restaurants like per one user like or has how many users like one restaurant
type Filter struct {
	RestaurantId int `json:"-" form:"restaurantId"`
	UserId       int `json:"-" form:"user_id"`
}
