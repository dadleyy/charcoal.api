package net

import "fmt"
import "strconv"
import "strings"
import "net/url"

import "github.com/jinzhu/gorm"
import "github.com/dadleyy/charcoal.api/util"

const BlueprintDefaultLimit = 100
const BlueprintMaxLimit = 500
const BlueprintFilterStart = "filter["
const BlueprintFilterEnd = "]"

type Blueprint struct {
	*gorm.DB
	values url.Values
}

func (print *Blueprint) Apply(out interface{}) (int, error) {
	var total int
	limit, page := BlueprintDefaultLimit, 0

	cursor := print.DB

	if i, err := strconv.Atoi(print.values.Get("limit")); err == nil {
		limit = util.MinInt(BlueprintMaxLimit, i)
	}

	if i, err := strconv.Atoi(print.values.Get("page")); err == nil {
		page = i
	}

	for key := range print.values {
		filterable := strings.HasPrefix(key, BlueprintFilterStart) && strings.HasSuffix(key, BlueprintFilterEnd)
		value := strings.SplitN(print.values.Get(key), "(", 2)

		if filterable == false || len(value) != 2 || strings.HasSuffix(value[1], ")") != true {
			continue
		}

		column := strings.TrimSuffix(strings.TrimPrefix(key, BlueprintFilterStart), BlueprintFilterEnd)
		operation, target := value[0], strings.TrimSuffix(value[1], ")")

		switch operation {
		case "in":
			values := strings.Split(target, ",")
			cursor = cursor.Where(fmt.Sprintf("%s in (?)", column), values)
		case "lk":
			query, search := fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%s%%", target)
			cursor = cursor.Where(query, search)
		case "eq":
			cursor = cursor.Where(fmt.Sprintf("%s = ?", column), target)
		case "lt":
			cursor = cursor.Where(fmt.Sprintf("%s < ?", column), target)
		case "gt":
			cursor = cursor.Where(fmt.Sprintf("%s > ?", column), target)
		}
	}

	// now that we've chained all our filters, execute the db query and return the error, if any
	if e := cursor.Limit(limit).Offset(page * limit).Find(out).Error; e != nil {
		return -1, e
	}

	// also make a `count()` request
	cursor.Model(out).Count(&total)

	return total, nil
}
