package models

type GameMembership struct {
	Common
	UserID uint `json:"user" gorm:"column:user_id"`
	GameID uint `json:"game" gorm:"column:game_id"`

	Game Game
	User User
}

func (membership GameMembership) TableName() string {
	return "game_memberships"
}
