package api

import "fmt"

import "github.com/sizethree/meritoss.api/api/db"

type Filter struct {
	Key string
	Operation string
	Value string
}

type Blueprint struct {
	Page int
	Limit int
	Filters []Filter
}

func (f *Filter) Reduce() string {
	return fmt.Sprintf("%s %s %s", f.Key, f.Operation, f.Value)
}

func (print *Blueprint) Apply(out interface{}, client *db.Client) (int, error) {
	var total int
	limit, offset := print.Limit, print.Limit * print.Page

	result := client.Begin().Limit(limit).Offset(offset)

	for _, filter := range print.Filters {
		result = result.Where(filter.Reduce())
	}

	e := result.Find(out).Count(&total).Error

	return total, e
}
