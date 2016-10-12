package models

type ClientToken struct {
	Common
	User uint `json:"user"`
	Token string `json:"token"`
	Client uint `json:"client"`
}

func (token ClientToken) Public() interface{} {
	return token
}
