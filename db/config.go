package db

import "fmt"

const DSN_STR = "%v:%v@tcp(%v:%v)/%v?parseTime=true&charset=utf8mb4,utf8"

type Config struct {
	Username string
	Password string
	Hostname string
	Database string
	Port     string
	Debug    bool
}

func (config *Config) String() string {
	return fmt.Sprintf(DSN_STR, config.Username, config.Password, config.Hostname, config.Port, config.Database)
}
