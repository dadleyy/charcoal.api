package net

import "fmt"
import "strconv"
import "strings"
import "net/url"

import "github.com/jinzhu/gorm"
import "github.com/gedex/inflector"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/util"

const BlueprintDefaultLimit = 100
const BlueprintMaxLimit = 500
const BlueprintMinLimit = 1
const BlueprintFilterStart = "filter["
const BlueprintFilterEnd = "]"

type Blueprint struct {
	*gorm.DB
	*log.Logger

	values url.Values
}

func (print *Blueprint) Limit() int {
	if i, err := strconv.Atoi(print.values.Get("limit")); err == nil {
		return util.MaxInt(util.MinInt(BlueprintMaxLimit, i), BlueprintMinLimit)
	}

	return BlueprintDefaultLimit
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

	scope := print.NewScope(out)
	table := scope.TableName()

	for key := range print.values {
		filterable := strings.HasPrefix(key, BlueprintFilterStart) && strings.HasSuffix(key, BlueprintFilterEnd)
		value := strings.SplitN(print.values.Get(key), "(", 2)

		if filterable == false || len(value) != 2 || strings.HasSuffix(value[1], ")") != true {
			continue
		}

		column := strings.TrimSuffix(strings.TrimPrefix(key, BlueprintFilterStart), BlueprintFilterEnd)
		operation, target := value[0], strings.TrimSuffix(value[1], ")")

		if bits := strings.Split(column, "."); len(bits) == 2 {
			print.Debugf("found an association query: %s - %s(%s)", column, operation, target)
			other, fk := inflector.Pluralize(bits[0]), fmt.Sprintf("%s_id", inflector.Singularize(bits[0]))
			join := fmt.Sprintf("JOIN %s ON %s.id = %s.%s", other, other, table, fk)
			column = fmt.Sprintf("%s.%s", other, bits[1])
			cursor = cursor.Joins(join)
		}

		full := fmt.Sprintf("%s.%s", table, column)

		switch operation {
		case "in":
			values := strings.Split(target, ",")
			cursor = cursor.Where(fmt.Sprintf("%s in (?)", full), values)
		case "lk":
			query, search := fmt.Sprintf("%s LIKE ?", full), fmt.Sprintf("%%%s%%", target)
			cursor = cursor.Where(query, search)
		case "eq":
			cursor = cursor.Where(fmt.Sprintf("%s = ?", full), target)
		case "lt":
			cursor = cursor.Where(fmt.Sprintf("%s < ?", full), target)
		case "gt":
			cursor = cursor.Where(fmt.Sprintf("%s > ?", full), target)
		}
	}

	direction := "ASC"

	if o := print.values.Get("sort_order"); o == "desc" || o == "desc" {
		direction = "DESC"
	}

	sort := fmt.Sprintf("id %s", direction)

	if on := print.values.Get("sort_on"); len(on) >= 1 {
		sort = fmt.Sprintf("%s %s", on, direction)
	}

	// now that we've chained all our filters, execute the db query and return the error, if any
	if e := cursor.Limit(limit).Offset(page * limit).Order(sort).Find(out).Error; e != nil {
		return -1, e
	}

	// also make a `count()` request
	cursor.Model(out).Count(&total)

	return total, nil
}
