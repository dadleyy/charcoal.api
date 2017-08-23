package models

type SystemEmailDomain struct {
	Common
	Domain string `json:"domain"`
}

func (domains SystemEmailDomain) TableName() string {
	return "system_email_domains"
}
