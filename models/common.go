package models

import "time"
import "github.com/jinzhu/gorm"
import "github.com/satori/go.uuid"

type Common struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Uuid      string     `json:"uuid"`
}

func (record Common) Public() interface{} {
	return record
}

func (record *Common) Identifier() string {
	return record.Uuid
}

func (record *Common) BeforeCreate(tx *gorm.DB) error {
	id := uuid.NewV4()
	record.Uuid = id.String()
	return nil
}
