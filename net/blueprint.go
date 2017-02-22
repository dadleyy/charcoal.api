package net

import "github.com/jinzhu/gorm"

type Blueprint struct {
	*gorm.DB

	limit   int
	page    int
	orderby string
	filters FilterList
}

type BlueprintFilter interface {
	Apply(*gorm.DB) *gorm.DB
	String() string
}

type FilterList []BlueprintFilter

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

	for _, filter := range print.filters {
		cursor = filter.Apply(cursor)
	}

	if e := cursor.Limit(limit).Offset(offset).Find(out).Error; e != nil {
		return -1, e
	}

	cursor.Model(out).Count(&total)

	return total, nil
}
