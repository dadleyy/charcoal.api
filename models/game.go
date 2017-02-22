package models

type Game struct {
	Common
	OwnerID uint `json:"owner_id" gorm:"column:owner_id"`

	GameMemberships []GameMembership `json:"-"`
	Owner           User             `json:"-"`
}
