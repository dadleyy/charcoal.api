package models

type GameMembership struct {
	Common
	UserID       uint  `json:"user_id" gorm:"column:user_id"`
	GameID       uint  `json:"game_id" gorm:"column:game_id"`
	EntryRoundID *uint `json:"entry_round_id"`

	Game Game `json:"-"`
	User User `json:"-"`
}

func (membership GameMembership) TableName() string {
	return "game_memberships"
}
