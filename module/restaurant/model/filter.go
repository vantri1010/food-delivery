package restaurantmodel

type Filter struct {
	UserId int   `json:"user_id,omitempty" form:"user_id"`
	Status []int `json:"-"`
}
