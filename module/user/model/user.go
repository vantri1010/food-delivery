package usermodel

import (
	"errors"
	"food-delivery/common"
)

const EntityName = "User"

type User struct {
	common.SQLModel `json:",inline"`
	Email           string        `json:"email" gorm:"column:email"`
	Password        string        `json:"-" gorm:"column:password"`
	Salt            string        `json:"-" gorm:"column:salt"`
	LastName        string        `json:"last_name" gorm:"column:last_name"`
	FistName        string        `json:"first_name" gorm:"column:first_name"`
	Phone           string        `json:"phone" gorm:"column:phone"`
	Role            string        `json:"role" gorm:"column:role"`
	Avatar          *common.Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

func (u *User) GetUserId() int {
	return u.Id
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRole() string {
	return u.Role
}

func (User) TableName() string {
	return "users"
}

func (u *User) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeUser)
}

type UserCreate struct {
	common.SQLModel `json:",inline"`
	Email           string        `json:"email" gorm:"column:email"`
	Password        string        `json:"password" gorm:"column:password"`
	Salt            string        `json:"-" gorm:"column:salt"`
	LastName        string        `json:"last_name" gorm:"column:last_name"`
	FistName        string        `json:"first_name" gorm:"column:first_name"`
	Role            string        `json:"-" gorm:"column:role"`
	Avatar          *common.Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

func (UserCreate) TableName() string {
	return User{}.TableName()
}

func (u *UserCreate) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeUser)
}

type UserLogin struct {
	Email    string `json:"email" form:"email" gorm:"column:email"`
	Password string `json:"password" form:"password" gorm:"column:password"`
}

func (UserLogin) TableName() string {
	return User{}.TableName()
}

//type Account struct {
//	AccesToken   *tokenprovider.Token `json:"acces_token"`
//	RefreshToken *tokenprovider.Token `json:"refresh_token"`
//}
//
//func NewAccount(at, rt *tokenprovider.Token) *Account {
//	return &Account{
//		AccesToken:   at,
//		RefreshToken: rt,
//	}
//}

var (
	ErrEmailOrPasswordInvalid = common.NewCustomError(
		errors.New("email or password invalid"),
		"email or password invalid",
		"ErrUsernameOrPasswordInvaid",
	)

	ErrEmailExisted = common.NewCustomError(
		errors.New("email has already existed"),
		"email has already existed",
		"ErrEmailExisted")
)
