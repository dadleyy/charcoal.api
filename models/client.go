package models

import "fmt"
import "github.com/jinzhu/gorm"
import "github.com/satori/go.uuid"

type Client struct {
	Common
	Name         string `json:"name"`
	ClientID     string `json:"-"`
	ClientSecret string `json:"-"`
	RedirectUri  string `json:"redirect_uri"`
	Description  string `json:"description"`
	System       bool   `json:"system"`
	Uuid         string `json:"uuid"`
}

func (client *Client) Url() string {
	return fmt.Sprintf("/clients?filter[id]=eq(%d)", client.ID)
}

func (client *Client) Identifier() string {
	return client.Uuid
}

func (client *Client) Type() string {
	return "application/vnd.miritos.client+json"
}

func (client *Client) BeforeCreate(tx *gorm.DB) error {
	id := uuid.NewV4()
	client.Uuid = id.String()
	return nil
}
