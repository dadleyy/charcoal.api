package context

import "time"
import "strconv"
import "strings"
import "net/http"
import "github.com/labstack/echo"

const DEFAULT_LIMIT int = 100

type Miritos struct {
	echo.Context
	DB *Database
	Errors ErrorList
	Meta MetaData
	Results ResultList
}

func (runtime *Miritos) Result(result Result) {
	runtime.Results = append(runtime.Results, result)
}

func (runtime *Miritos) SetMeta(key string, value interface{}) {
	runtime.Meta[key] = value
}

func (runtime *Miritos) Error(err error) {
	runtime.Errors = append(runtime.Errors, err)
}

func (runtime *Miritos) Blueprint() Blueprint {
	result := Blueprint{Limit: DEFAULT_LIMIT, Page: 0}

	limit := runtime.QueryParam("limit")
	page := runtime.QueryParam("page")

	result.OrderBy = runtime.QueryParam("orderby")

	if value, err := strconv.Atoi(page); err == nil {
		result.Page = value
	}

	if value, err := strconv.Atoi(limit); err == nil {
		result.Limit = value
	}

	for key, values := range runtime.QueryParams() {
		filterable := strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]")

		if filterable == false || len(values) > 1 {
			continue
		}

		value := values[0]

		if err := result.Filter(key, value); err != nil {
			runtime.Error(err)
		}
	}

	return result
}

func (runtime *Miritos) Finish() error {
	runtime.Meta["time"] = time.Now()
	runtime.DB.Close()

	if ecount := len(runtime.Errors); ecount >= 1 {
		elist := make([]string, ecount)

		for indx, e := range runtime.Errors {
			elist[indx] = e.Error()
		}

		response := struct {
			Meta MetaData `json:"meta"`
			Status string `json:"status"`
			Errors []string `json:"errors"`
		}{runtime.Meta, "FAILED", elist}

		return runtime.JSON(http.StatusBadRequest, response)
	}

	if bp := runtime.Blueprint(); len(bp.Filters) >= 1 {
		filters := make([]string, len(bp.Filters) - 1)

		for _, f := range bp.Filters {
			filters = append(filters, f.String())
		}

		runtime.Meta["filters"] = filters
	}

	response := struct {
		Meta MetaData `json:"meta"`
		Status string `json:"status"`
		Results ResultList `json:"results"`
	}{runtime.Meta, "SUCCESS", runtime.Results}

	return runtime.JSON(http.StatusOK, response)
}