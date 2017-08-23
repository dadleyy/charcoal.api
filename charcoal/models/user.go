package models

import "fmt"
import "github.com/jinzhu/gorm"
import "github.com/docker/docker/pkg/namesgenerator"

type User struct {
	Common
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Username string `json:"password,omitempty"`

	GameMemberships []GameMembership `json:"-"`

	Games []Game `json:"-" gorm:"many2many:game_memberships;ForeignKey:id;AssociationForeignKey:id"`
}

func (user *User) Public() interface{} {
	out := struct {
		Common
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}{user.Common, user.Name, user.Email, user.Username}
	return out
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user == nil {
		return fmt.Errorf("received nil reference")
	}

	if user.Username != "" {
		return user.Common.BeforeCreate(tx)
	}

	user.Username = namesgenerator.GetRandomName(0)
	return user.Common.BeforeCreate(tx)
}

func (user *User) URL() string {
	return fmt.Sprintf("/users?filter[id]=eq(%d)", user.ID)
}

func (user *User) Type() string {
	return "application/vnd.miritos.user+json"
}
