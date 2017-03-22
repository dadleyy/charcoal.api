package models

type ClientToken struct {
	Common
	UserID   uint   `json:"user_id"`
	ClientID uint   `json:"client_id"`
	Token    string `json:"token"`
}
