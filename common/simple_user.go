package common

// this model using for preload user when get restaurants.
// to avoid getting sensitive information such as emails and phone numbers
// thus, we only declare what we need ( not re-use model.user )
type SimpleUser struct {
	SQLModel  `json:",inline"`
	LastName  string `json:"last_name" gorm:"column:last_name"`
	FirstName string `json:"first_name" gorm:"column:first_name"`
	Role      string `json:"role" gorm:"column:role"`
	Avatar    *Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

func (SimpleUser) TableName() string {
	return "users"
}

func (u *SimpleUser) Mask(isAdmin bool) {
	u.GenUID(DbTypeUser)
}
