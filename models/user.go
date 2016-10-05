package models

import "fmt"

type User struct {
	Common
	Name *string `json:"name",omitempty`
	Email *string `json:"email",omitempty`
	Password *string `json:"password",omitempty`
}

func (user *User) Public() interface{} {
	out := struct {
		Common
		Name string `json:"name"`
		Email string `json:"email"`
	}{user.Common, *user.Name, *user.Email}
	return out
}

func (user *User) Url() string {
	return fmt.Sprintf("http://example.com/users?filters[id]=eq(%d)", user.ID)
}

func (user *User) Type() string {
	return "application/vnd.miritos.user+json"
}
