package models

type UserRoleMapping struct {
	Common
	User uint `json:"user"`
	Role uint `json:"role"`
}

func (mapping UserRoleMapping) Public() interface{} {
	return mapping
}
