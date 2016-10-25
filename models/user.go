package models

import "os"
import "fmt"

type User struct {
	Common
	Name     *string `json:"name",omitempty`
	Email    *string `json:"email",omitempty`
	Password *string `json:"password",omitempty`
}

func (user User) Public() interface{} {
	out := struct {
		Common
		Name  string `json:"name"`
		Email string `json:"email"`
	}{user.Common, *user.Name, *user.Email}
	return out
}

func (user User) Url() string {
	root := os.Getenv("API_HOME")
	return fmt.Sprintf("%s/users?filter[id]=eq(%d)", root, user.ID)
}

func (user User) Type() string {
	return "application/vnd.miritos.user+json"
}
