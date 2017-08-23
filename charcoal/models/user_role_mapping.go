package models

type UserRoleMapping struct {
	Common
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
}
