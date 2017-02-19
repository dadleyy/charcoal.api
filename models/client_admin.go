package models

type ClientAdmin struct {
	Common
	User   uint `json:"user"`
	Client uint `json:"client"`
}
