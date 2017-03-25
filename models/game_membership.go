package models

type GameMembership struct {
	Common
	UserID uint   `json:"user_id" gorm:"column:user_id"`
	GameID uint   `json:"game_id" gorm:"column:game_id"`
	Status string `json:"status"`

	Game Game `json:"-"`
	User User `json:"-"`
}

func (membership GameMembership) TableName() string {
	return "game_memberships"
}
