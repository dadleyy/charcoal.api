package models

type User struct {
	Common
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"-"`
}
