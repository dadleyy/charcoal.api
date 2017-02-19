package models

type Game struct {
	Common
	OwnerID uint `json:"owner" gorm:"column:owner_id"`

	Owner User `json:"-"`
}
