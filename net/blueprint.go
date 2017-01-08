package net

import "fmt"
import "strings"
import "github.com/sizethree/miritos.api/db"

type Blueprint struct {
	limit   int
	page    int
	orderby string
	filters FilterList
}

type BlueprintFilter interface {
	Apply(*db.Connection) *db.Connection
	String() string
}

type FilterList []BlueprintFilter

func (print *Blueprint) Limit() int {
	return print.limit
}

func (print *Blueprint) Page() int {
	return print.page
}

func (print *Blueprint) Filter(key string, opstr string) error {
	// extract the field from the filter string key
	field := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(key, "filter["), "]"))

	if len(field) < 2 {
		return fmt.Errorf("INVALID_FIELD - [%s]", field)
	}

	// all filters should look like a function call: `op(value)`
	parts := strings.Split(opstr, "(")

	if len(parts) != 2 || strings.HasSuffix(parts[1], ")") != true {
		return fmt.Errorf("BAD_OPERATION")
	}

	value := strings.TrimSuffix(parts[1], ")")

	// make sure we have a valid value
	if lv := len(value); lv < 1 {
		return fmt.Errorf("BAD_OPERATION_PARTS")
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
	case "eq":
		filter = &sizeOp{field, value, "="}
	case "in":
		filter = &inOp{field, value}
	}

	if filter.String() != "" {
		print.filters = append(print.filters, filter)
	}

	return nil
}

func (print *Blueprint) Apply(out interface{}, client *db.Connection) (int, error) {
	var total int
	limit, offset := print.limit, print.limit*print.page

	result := client

	for _, filter := range print.filters {
		result = filter.Apply(result)
	}

	e := result.Limit(limit).Offset(offset).Find(out).Error
	result.Model(out).Count(&total)

	return total, e
}

type inOp struct {
	field string
	value string
}

func (op *inOp) Apply(client *db.Connection) *db.Connection {
	ins := strings.Split(op.value, ",")
	query := fmt.Sprintf("%s in (?)", op.field)
	return &db.Connection{client.Where(query, []string(ins))}
}

func (op *inOp) String() string {
	return fmt.Sprintf("%s in (%s)", op.field, op.value)
}

// sizeOp
//
// given a field, value and an operator, this operation simply uses the database
// client's `Where` function with appropriate string format.
type sizeOp struct {
	field    string
	value    string
	operator string
}

func (op *sizeOp) Apply(client *db.Connection) *db.Connection {
	clause := fmt.Sprintf("%s %s ?", op.field, op.operator)
	return &db.Connection{client.Where(clause, op.value)}
}

func (op *sizeOp) String() string {
	return fmt.Sprintf("%s %s %s", op.field, op.operator, op.value)
}

// noOp
//
// describes a blueprint filter operation that is not understood by the system.
type noOp struct{}

func (op *noOp) String() string {
	return ""
}

func (op *noOp) Apply(client *db.Connection) *db.Connection {
	return client
}
