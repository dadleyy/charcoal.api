package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name string
	Email string
	Password string
}

func (u User) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id": u.ID,
		"name": u.Name,
		"email": u.Email,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}
