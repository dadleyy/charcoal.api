package models

type SystemSettings struct {
	Common
	RestrictedEmailDomains bool `json:"restricted_email_domains"`
}

func (settings SystemSettings) TableName() string {
	return "system_settings"
}
