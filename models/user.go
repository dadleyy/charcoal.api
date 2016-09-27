package models

import "time"

type User struct {
	Common
	Name string
	Email string
	Password string
}

type UserMarshal struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (user *User) Marshal() interface{} {
	return &UserMarshal{user.ID, user.Name, user.Email, user.UpdatedAt, user.CreatedAt}
}
