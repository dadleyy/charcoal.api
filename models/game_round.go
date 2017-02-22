package models

import "database/sql"

type GameRound struct {
	Common
	GameID      uint          `json:"game_id" gorm:"column:game_id"`
	PresidentID sql.NullInt64 `json:"president" gorm:"column:president_id"`

	Game Game `json:"-"`
}

func (r *GameRound) Public() interface{} {
	var pres interface{} = nil

	if r.PresidentID.Valid {
		pres = r.PresidentID.Int64
	}

	return struct {
		Common
		President interface{} `json:"president_id"`
		Game      interface{} `json:"game_id"`
	}{r.Common, pres, r.GameID}
}
