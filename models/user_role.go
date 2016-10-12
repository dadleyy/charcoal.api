package models

type UserRole struct {
	Common
	Label string `json:"label"`
	Description string `json:"description"`
}

func (role UserRole) Public() interface{} {
	return role
}
