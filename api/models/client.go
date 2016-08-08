package models

import "github.com/jinzhu/gorm"

type Client struct {
	gorm.Model
	Name string
	ClientID string
	ClientSecret string
}

func (client Client) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id": client.ID,
		"name": client.Name,
		"created_at": client.CreatedAt,
		"updated_at": client.UpdatedAt,
	}
}
