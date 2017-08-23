package models

type GameMembership struct {
	Common
	UserID uint   `json:"user_id" gorm:"column:user_id"`
	GameID uint   `json:"game_id" gorm:"column:game_id"`
	Status string `json:"status"`

	Assholeships     int `json:"assoleships"`
	Presidencies     int `json:"presidencies"`
	VicePresidencies int `json:"vice_presidencies"`

	Game Game `json:"-"`
	User User `json:"-"`
}

func (membership GameMembership) TableName() string {
	return "game_memberships"
}
