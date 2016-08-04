package middleware

import "math"
import "strings"
import "strconv"

import "github.com/golang/glog"
import "github.com/kataras/iris"

const DEFAULT_PAGE = 0
const DEFAULT_LIMIT = 100
const MAX_LIMIT = 300

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

// Blueprints
// 
// given an iris context, this function will register a user value `blueprints` into 
// it that is a `Blueprint` struct defined above. these are useful in resource lookup 
// routes (e.g GET /users)
func Blueprints(ctx *iris.Context) {
	glog.Info("parsing blueprints")

	blueprints := Blueprint{Page: DEFAULT_PAGE, Limit: DEFAULT_LIMIT}

	params := ctx.URLParams()

	for key, value := range params {
		glog.Infof("found param[%s] value[%s]\n", key, value)

		switch strings.ToLower(key) {
		case "max":
			ival, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			blueprints.Limit = ival
		case "page":
			ival, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			blueprints.Page = ival
		default:
			continue;
		}

	}

	blueprints.Limit = int(math.Min(float64(MAX_LIMIT), float64(blueprints.Limit)))

	ctx.Set("blueprint", blueprints)

	ctx.Next()
}
