package context

import "github.com/jinzhu/gorm"

type Database struct {
	*gorm.DB
}

func (client *Database) Where(clause, values string) *Database {
	result := client.DB.Where(clause, values)
	return &Database{result}
}
