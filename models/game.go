package models

import "fmt"

const GameDefaultStatus = "ACTIVE"

type Game struct {
	Common
	OwnerID    uint   `json:"owner_id" gorm:"column:owner_id"`
	Status     string `json:"status"`
	Population int    `json:"popuplation"`

	GameMemberships []GameMembership `json:"-"`
	Owner           User             `json:"-"`
}

func (g *Game) Url() string {
	return fmt.Sprintf("/games?filter[id]=eq(%d)", g.ID)
}

func (g *Game) Type() string {
	return "application/vnd.charcoal.game+json"
}
