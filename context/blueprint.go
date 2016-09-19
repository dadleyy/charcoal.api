package context

import "fmt"
import "errors"
import "strings"

type Blueprint struct {
	Limit int
	Page int
	OrderBy string
	Filters FilterList
}

type BlueprintFilter interface {
	Apply(*Database) *Database
	String() string
}

type FilterList []BlueprintFilter

func (print *Blueprint) Filter(key string, opstr string) error {
	// extract the field from the filter string key
	field := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(key, "filter["), "]"))

	if len(field) < 2 {
		return errors.New(fmt.Sprintf("INVALID_FIELD - [%s]", field))
	}

	// all filters should look like a function call: `op(value)`
	parts := strings.Split(opstr, "(")

	if len(parts) != 2 || strings.HasSuffix(parts[1], ")") != true {
		return errors.New("BAD_OPERATION")
	}

	value := strings.TrimSuffix(parts[1], ")")

	// make sure we have a valid value
	if lv := len(value); lv < 1 {
		return errors.New("BAD_OPERATION_PARTS")
	}

	var filter BlueprintFilter = &noOp{}

	switch parts[0] {
	case "gt":
		filter = &sizeOp{field, value, ">"}
	case "gte":
		filter = &sizeOp{field, value, ">="}
	case "lt":
		filter = &sizeOp{field, value, "<"}
	case "lte":
		filter = &sizeOp{field, value, "<="}
	}

	if filter.String() != "" {
		print.Filters = append(print.Filters, filter)
	}

	return nil
}


func (print *Blueprint) Apply(out interface{}, client *Database) (int, error) {
	var total int
	limit, offset := print.Limit, print.Limit * print.Page

	result := &Database{client.Begin().Limit(limit).Offset(offset)}

	for _, filter := range print.Filters {
		result = filter.Apply(result)
	}

	e := result.Find(out).Count(&total).Error

	return total, e
}

type sizeOp struct {
	field string
	value string
	operator string
}

func (op *sizeOp) Apply(client *Database) *Database {
	clause := fmt.Sprintf("%s %s ?", op.field, op.operator)
	return client.Where(clause, op.value)
}

func (op *sizeOp) String() string {
	return fmt.Sprintf("%s %s %s", op.field, op.operator, op.value)
}

type noOp struct {}

func (op *noOp) String() string {
	return ""
}

func (op *noOp) Apply(client *Database) *Database {
	return client
}


