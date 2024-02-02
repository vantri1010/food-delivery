package common

import "log"

func AppRecover() {
	if err := recover(); err != nil {
		log.Println("recover error: ", err)
	}
}
