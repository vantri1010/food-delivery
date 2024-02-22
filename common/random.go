package common

import (
	"math/rand"
	"time"
)

// function create a random string of characters only (uppercase or lowercase), no numbers
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randSequence(n int) string {
	b := make([]byte, n)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := range b {
		b[i] = letters[r1.Int63()%int64(len(letters))]
	}

	return string(b)
}

func GenSalt(lengh int) string {
	if lengh < 0 {
		lengh = 50
	}

	return randSequence(lengh)
}
