package models

type GameRound struct {
	Common
	GameID uint `json:"game_id" gorm:"column:game_id"`

	PresidentID     *int64 `json:"president_id" gorm:"column:president_id"`
	VicePresidentID *int64 `json:"vice_president_id" gorm:"column:vice_president_id"`
	AssholeID       *int64 `json:"asshole_id" gorm:"column:asshole_id"`

	Game Game `json:"-"`
}

func (r *GameRound) Public() interface{} {
	return *r
}
