package models

import "github.com/jinzhu/gorm"

type PositionHistory struct {
	gorm.Model
	Before int
	After int
	User uint `gorm:"ForeignKey:users.ID;"`
	Proposal uint `gorm:"ForeignKey:proposals.ID;"`
}

func (pos PositionHistory) TableName() string {
	return "position_history"
}

func (pos PositionHistory) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id": pos.ID,
		"before": pos.Before,
		"after": pos.After,
		"user": pos.User,
		"proposal": pos.Proposal,
		"created_at": pos.CreatedAt,
		"updated_at": pos.UpdatedAt,
	}
}
