package models

import "time"
import "github.com/jinzhu/gorm"
import "github.com/satori/go.uuid"

// Common is a struct definition that contains basic fields shared across all model types.
type Common struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UUID      string     `json:"uuid"`
}

// Public returns the record w/ any private fields omitted.
func (record Common) Public() interface{} {
	return record
}

// Identifier returns unique id associated w/ the record.
func (record *Common) Identifier() string {
	return record.UUID
}

// BeforeCreate is a hook used by GORM; takes care of auto-uuid generation.
func (record *Common) BeforeCreate(tx *gorm.DB) error {
	id := uuid.NewV4()
	record.UUID = id.String()
	return nil
}
