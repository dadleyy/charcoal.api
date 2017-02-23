package models

type Game struct {
	Common
	OwnerID uint   `json:"owner_id" gorm:"column:owner_id"`
	Status  string `json:"status"`

	GameMemberships []GameMembership `json:"-"`
	Owner           User             `json:"-"`
}
