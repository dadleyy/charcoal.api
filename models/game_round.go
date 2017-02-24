package models

import "database/sql"

type GameRound struct {
	Common
	GameID          uint          `json:"game_id" gorm:"column:game_id"`
	PresidentID     sql.NullInt64 `json:"president_id" gorm:"column:president_id"`
	VicePresidentID sql.NullInt64 `json:"vice_president_id" gorm:"column:vice_president_id"`
	AssholeID       sql.NullInt64 `json:"asshole_id" gorm:"column:asshole_id"`

	Game Game `json:"-"`
}

func (r *GameRound) Public() interface{} {
	var pres, vp, ass interface{} = nil, nil, nil

	if r.PresidentID.Valid {
		pres = r.PresidentID.Int64
	}

	if r.AssholeID.Valid {
		ass = r.AssholeID.Int64
	}

	if r.VicePresidentID.Valid {
		vp = r.VicePresidentID.Int64
	}

	return struct {
		Common
		Vice interface{} `json:"vice_president_id"`
		Pres interface{} `json:"president_id"`
		Game interface{} `json:"game_id"`
		Ass  interface{} `json:"asshole_id"`
	}{r.Common, vp, pres, r.GameID, ass}
}
