package models

type User struct {
	Common
	Name *string `json:"name",omitempty`
	Email *string `json:"email",omitempty`
	Password *string `json:"password",omitempty`
}

func (user *User) Public() interface{} {
	out := struct {
		Common
		Name string `json:"name"`
		Email string `json:"email"`
	}{user.Common, *user.Name, *user.Email}
	return out
}
