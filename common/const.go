package common

import "log"

const (
	DbTypeRestaurant = 1
	DbTypeUser       = 2
)

const (
	CurrentUser = "user"
)

func AppRecover() {
	if err := recover(); err != nil {
		log.Println("recover error: ", err)
	}
}
