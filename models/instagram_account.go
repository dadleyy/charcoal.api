package models

import "database/sql"
import "github.com/jinzhu/gorm"
import "github.com/satori/go.uuid"

type InstagramAccount struct {
	Common
	InstagramID string        `json:"instagram_id" gorm:"column:instagram_id"`
	User        sql.NullInt64 `json:"user"`
	Username    string        `json:"username"`
	Uuid        string        `json:"uuid"`
}

func (account *InstagramAccount) BeforeCreate(tx *gorm.DB) error {
	id := uuid.NewV4()
	account.Uuid = id.String()
	return nil
}

func (account *InstagramAccount) Public() interface{} {
	var user interface{} = nil

	if account.User.Valid {
		user = account.User.Int64
	}

	return struct {
		InstagramAccount
		User interface{} `json:"user"`
	}{*account, user}
}
