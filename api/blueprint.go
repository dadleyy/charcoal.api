package api

import "fmt"
import "github.com/golang/glog"
import "github.com/jinzhu/gorm"

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

func (print *Blueprint) Apply(runtime *Runtime) *gorm.DB {
	limit, offset := print.Limit, print.Limit * print.Page

	result := runtime.DB.Begin().Limit(limit).Offset(offset)

	for _, filter := range print.Filters {
		glog.Infof("adding filter \"%s\"\n", filter.Reduce())
		result = result.Where(filter.Reduce())
	}

	return result
}
