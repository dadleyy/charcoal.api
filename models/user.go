package models

import "time"
import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name string
	Email string
	Password string
}

type UserMarshal struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	UpdatedAt time.Time `json:"updated"`
	CreatedAt time.Time `json:"created"`
}

func (user *User) Marshal() interface{} {
	return &UserMarshal{user.ID, user.Name, user.Email, user.UpdatedAt, user.CreatedAt}
}
