package models

import "github.com/jinzhu/gorm"

type Proposal struct {
	gorm.Model
	Summary string
	Content string
	Author uint `gorm:"ForeignKey:users.ID;"`
}

func (p Proposal) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id": p.ID,
		"summary": p.Summary,
		"content": p.Content,
		"author": p.Author,
		"created_at": p.CreatedAt,
		"updated_at": p.UpdatedAt,
	}
}

