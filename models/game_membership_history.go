package models

type GameMembershipHistory struct {
	Common
	UserID       uint  `json:"user_id" gorm:"column:user_id"`
	GameID       uint  `json:"game_id" gorm:"column:game_id"`
	EntryRoundID *uint `json:"entry_round_id"`
	ExitRoundID  *uint `json:"exit_round_id"`
}

func (h GameMembershipHistory) TableName() string {
	return "game_membership_history"
}
