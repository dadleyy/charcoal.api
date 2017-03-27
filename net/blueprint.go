package net

import "fmt"
import "strconv"
import "strings"
import "net/url"

import "github.com/jinzhu/gorm"
import "github.com/gedex/inflector"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/util"
import "github.com/dadleyy/charcoal.api/defs"

type Blueprint struct {
	*gorm.DB
	*log.Logger

	values url.Values
}

type BlueprintForeignReference struct {
	reference string
	source    string
}

func (r *BlueprintForeignReference) JoinString() string {
	bits := strings.Split(r.reference, ".")
	table := inflector.Pluralize(bits[0])
	fk := fmt.Sprintf("%s_id", inflector.Singularize(table))
	return fmt.Sprintf("JOIN %s on %s.id = %s.%s", table, table, r.source, fk)
}

func (r *BlueprintForeignReference) WhereField() string {
	bits := strings.Split(r.reference, ".")
	table := inflector.Pluralize(bits[0])
	return fmt.Sprintf("%s.%s", table, bits[1])
}

func (print *Blueprint) Limit() int {
	if i, err := strconv.Atoi(print.values.Get("limit")); err == nil {
		return util.MaxInt(util.MinInt(defs.BlueprintMaxLimit, i), defs.BlueprintMinLimit)
	}

	return defs.BlueprintDefaultLimit
}

func (print *Blueprint) Apply(out interface{}) (int, error) {
	var total int
	limit, page := defs.BlueprintDefaultLimit, 0

	cursor := print.DB

	if i, err := strconv.Atoi(print.values.Get("limit")); err == nil {
		limit = util.MinInt(defs.BlueprintMaxLimit, i)
	}

	if i, err := strconv.Atoi(print.values.Get("page")); err == nil {
		page = i
	}

	scope := print.NewScope(out)
	table := scope.TableName()

	for key := range print.values {
		filterable := strings.HasPrefix(key, defs.BlueprintFilterStart) && strings.HasSuffix(key, defs.BlueprintFilterEnd)
		value := strings.SplitN(print.values.Get(key), "(", 2)

		if filterable == false || len(value) != 2 || strings.HasSuffix(value[1], ")") != true {
			continue
		}

		column := strings.TrimSuffix(strings.TrimPrefix(key, defs.BlueprintFilterStart), defs.BlueprintFilterEnd)
		operation, target := value[0], strings.TrimSuffix(value[1], ")")
		full := fmt.Sprintf("%s.%s", table, column)

		if bits := strings.Split(column, "."); len(bits) == 2 {
			print.Debugf("found an association query: %s - %s(%s)", column, operation, target)
			reference := BlueprintForeignReference{column, table}

			// move the cursor into a join + change the where clause to be our referenced where
			cursor = cursor.Joins(reference.JoinString())
			full = reference.WhereField()
		}

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
