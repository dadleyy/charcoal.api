package models

import "github.com/jinzhu/gorm"

type ClientToken struct {
	gorm.Model
	Token string
	User uint `gorm:"ForeignKey:users.ID;"`
	Client uint `gorm:"ForeignKey:clients.ID;"`
}

func (token ClientToken) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id": token.ID,
		"token": token.Token,
		"user": token.User,
		"client": token.Client,
		"created_at": token.CreatedAt,
		"updated_at": token.UpdatedAt,
	}
}

