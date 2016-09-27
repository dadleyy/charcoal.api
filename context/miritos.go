package context

import "time"
import "strconv"
import "strings"
import "net/http"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"

const DEFAULT_LIMIT int = 100

type Miritos struct {
	echo.Context
	DB *Database
	Errors ErrorList
	Meta MetaData
	Results ResultList
	FS FileSaver
}

func (runtime *Miritos) Body() (Body, error) {
	body := make(Body)
	err := runtime.Bind(&body)
	return body, err
}

func (runtime *Miritos) Result(result Result) {
	runtime.Results = append(runtime.Results, result)
}

func (runtime *Miritos) SetMeta(key string, value interface{}) {
	runtime.Meta[key] = value
}

func (runtime *Miritos) ErrorOut(err error) error {
	runtime.Errors = append(runtime.Errors, err)
	return nil
}

func (runtime *Miritos) PersistFile(target File) (models.File, error) {
	temp, err := runtime.FS.Upload(target)

	if err != nil {
		return temp, err
	}

	if err := runtime.DB.Create(&temp).Error; err != nil {
		return temp, err
	}

	return temp, err
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

	results := make([]interface{}, len(runtime.Results))

	for i, result := range runtime.Results {
		results[i] = result.Marshal()
	}

	response := struct {
		Meta MetaData `json:"meta"`
		Status string `json:"status"`
		Results []interface{} `json:"results"`
	}{runtime.Meta, "SUCCESS", results}

	return runtime.JSON(http.StatusOK, response)
}
