package net

import "github.com/jinzhu/gorm"

type Blueprint struct {
	*gorm.DB

	limit   int
	page    int
	orderby string
}

func (print *Blueprint) Limit() int {
	return print.limit
}

func (print *Blueprint) Page() int {
	return print.page
}

func (print *Blueprint) Apply(out interface{}) (int, error) {
	var total int
	limit, offset := print.limit, print.limit*print.page
	cursor := print.DB

	if e := cursor.Limit(limit).Offset(offset).Find(out).Error; e != nil {
		return -1, e
	}

	cursor.Model(out).Count(&total)

	return total, nil
}
