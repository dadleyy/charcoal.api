package models

type ClientAdmin struct {
	Common
	UserID   uint `json:"user_id"`
	ClientID uint `json:"client_id"`
}
