package models

import "fmt"
import "github.com/jinzhu/gorm"
import "github.com/satori/go.uuid"

type User struct {
	Common
	Name     *string `json:"name",omitempty`
	Email    *string `json:"email",omitempty`
	Password *string `json:"password",omitempty`
	Uuid     *string `json:"uuid",omitempty`
}

func (user *User) Public() interface{} {
	out := struct {
		Common
		Name  string `json:"name"`
		Email string `json:"email"`
	}{user.Common, *user.Name, *user.Email}
	return out
}

func (user *User) Identifier() string {
	if user.Uuid == nil {
		return ""
	}

	return *user.Uuid
}

func (user *User) Url() string {
	return fmt.Sprintf("/users?filter[id]=eq(%d)", user.ID)
}

func (user *User) Type() string {
	return "application/vnd.miritos.user+json"
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	id := uuid.NewV4().String()
	user.Uuid = &id
	return nil
}
