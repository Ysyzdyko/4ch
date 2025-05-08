package core

import (
	"errors"
	"math/rand"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (user *User) GenUserID(userCount, maxID int) (int, error) {
	if maxID == 0 {
		return 0, errors.New("maxID is zero")
	}

	id := userCount + 1
	if userCount > maxID {
		id = rand.Intn(maxID)
	}

	return id, nil
}
