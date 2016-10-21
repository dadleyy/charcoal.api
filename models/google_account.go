package models

type GoogleAccount struct {
	Common
	GoogleID    string `json:"google_id" gorm:"column:google_id"`
	User        uint   `json:"user"`
	AccessToken string `json:"-"`
	Email       string `json:"email"`
	Name        string `json:"name"`
}
