package models

import "fmt"
import "database/sql"
import "github.com/jinzhu/gorm"
import "github.com/satori/go.uuid"

type serializedInstagramPhoto struct {
	InstagramPhoto
	Client interface{} `json:"client"`
}

type InstagramPhoto struct {
	Common
	Owner       string        `json:"owner"`
	InstagramID string        `json:"instagram_id"`
	Caption     string        `json:"caption"`
	Photo       uint          `json:"photo"`
	Client      sql.NullInt64 `json:"client"`
	Uuid        string        `json:"uuid"`
}

func (photo *InstagramPhoto) Identifier() string {
	return photo.Uuid
}

func (photo *InstagramPhoto) Url() string {
	return fmt.Sprintf("/instagram?filter[id]=eq(%d)", photo.ID)
}

func (photo *InstagramPhoto) Type() string {
	return "application/vnd.miritos.instagram-photo+json"
}

func (photo *InstagramPhoto) BeforeCreate(tx *gorm.DB) error {
	id := uuid.NewV4()
	photo.Uuid = id.String()
	return nil
}

func (photo *InstagramPhoto) Public() interface{} {
	var author interface{}

	author = nil
	if photo.Client.Valid {
		author = photo.Client.Int64
	}

	result := serializedInstagramPhoto{*photo, author}
	return result
}
