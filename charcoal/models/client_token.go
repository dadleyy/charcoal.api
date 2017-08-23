package models

// ClientToken records represent user-given permission for an application client to communicate on their behalf.
type ClientToken struct {
	Common
	UserID   uint   `json:"user_id"`
	ClientID uint   `json:"client_id"`
	Token    string `json:"token"`
}
