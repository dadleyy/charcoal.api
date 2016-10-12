package models

type ClientAdmin struct {
	Common
	User uint `json:"user"`
	Client uint `json:"client"`
}

func (admin ClientAdmin) Public() interface{} {
	return admin
}
