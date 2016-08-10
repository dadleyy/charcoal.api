package middleware

import "fmt"
import "math"
import "errors"
import "strings"
import "strconv"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/db"

const (
	DEFAULT_PAGE = 0
	DEFAULT_LIMIT = 100
	MAX_LIMIT = 300
)

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

// intOr
// 
// returns either the value parsed as an integer or the user-provided "backup" value
func intOr(value string, backup int) int {
	ival, err := strconv.Atoi(value)

	if err != nil {
		return backup
	}

	return ival
}

func parseFilter(key string, value string) (Filter, error) {
	fieldkey := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(key, "filter["), "]"))

	if len(fieldkey) < 1 {
		return Filter{}, errors.New("BAD_FILTER_KEY")
	}

	// to find out the value and the operation, we need to split the string by "(" since 
	// these are called out in the form: op(val)
	parts := strings.Split(value, "(")

	if len(parts) != 2 || strings.HasSuffix(parts[1], ")") != true {
		return Filter{}, errors.New("BAD_OPERATION")
	}

	fieldval, fieldop := strings.TrimSuffix(parts[1], ")"), ""

	switch parts[0] {
	case "gt":
		fieldop = ">"
	case "lt":
		fieldop = "<"
	case "eq":
		fieldop = "="
	case "neq":
		fieldop = "!="
	case "lk":
		fieldop = "like"
	case "gte":
		fieldop = ">="
	case "lte":
		fieldop = "<="
	case "in":
		fieldop = "in"
		fieldval = fmt.Sprintf("(%s)", fieldval)
	}

	// make sure that our filter's value exists and our operation is either 2 or 3 characters
	// long; these are the lengths of our operations
	if lv, lo := len(fieldval), len(fieldop); lv < 1 || lo < 1 {
		return Filter{}, errors.New("BAD_OPERATION_PARTS")
	}

	return Filter{fieldkey, fieldop, fieldval}, nil
}

// InjectBlueprint
// 
// given an iris context, this function will register a user value `blueprints` into 
// it that is a `Blueprint` struct defined above. these are useful in resource lookup 
// routes (e.g GET /users)
func InjectBlueprint(ctx *iris.Context) {
	blueprint := Blueprint{Page: DEFAULT_PAGE, Limit: DEFAULT_LIMIT}

	for key, value := range ctx.URLParams() {
		cleankey := strings.TrimSpace(strings.ToLower(key))
		cleanval := strings.TrimSpace(value)


		// based on the trimmed/lowercasd key of this parameter we will either want to set
		// the blueprint's limit, page, items in it's Filters slice or nothing at all.
		switch {
		case cleankey == "max":
			blueprint.Limit = intOr(cleanval, DEFAULT_LIMIT)
		case cleankey == "page":
			blueprint.Page = intOr(cleanval, DEFAULT_PAGE)
		case strings.HasPrefix(cleankey, "filter[") && strings.HasSuffix(cleankey, "]"):
			if fil, err := parseFilter(cleankey, cleanval); err == nil {
				blueprint.Filters = append(blueprint.Filters, fil)
			}
		default:
			continue;
		}

	}

	blueprint.Limit = int(math.Min(float64(MAX_LIMIT), float64(blueprint.Limit)))

	ctx.Set("blueprint", &blueprint)

	ctx.Next()
}
