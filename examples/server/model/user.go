package model

import (
	"fmt"
	"time"
)

var _ Model = (*User)(nil)

type (
	// User represents a user in the system
	User struct {
		ID        uint      `json:"id"`
		Name      string    `json:"name"`
		Birthdate time.Time `json:"birthdate"`
	}
)

func (u *User) SetID(id uint) {
	u.ID = id
}

func (u *User) GetID() uint {
	return u.ID
}

func (u *User) String() string {
	return fmt.Sprintf(
		"User{ID: %d, Name: %s, Birthdate: %s}",
		u.ID,
		u.Name,
		u.Birthdate,
	)
}
