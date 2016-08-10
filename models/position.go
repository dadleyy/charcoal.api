package models

import "github.com/jinzhu/gorm"

type Position struct {
	gorm.Model
	Location int
	User uint `gorm:"ForeignKey:users.ID;"`
	Proposal uint `gorm:"ForeignKey:proposals.ID;"`
}

func (pos Position) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id": pos.ID,
		"location": pos.Location,
		"user": pos.User,
		"proposal": pos.Proposal,
		"created_at": pos.CreatedAt,
		"updated_at": pos.UpdatedAt,
	}
}
