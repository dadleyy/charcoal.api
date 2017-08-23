package models

// ClientAdmin records are associations b/w clients and users that reprent users w/ the ability to manage the client.
type ClientAdmin struct {
	Common
	UserID   uint `json:"user_id"`
	ClientID uint `json:"client_id"`
}
