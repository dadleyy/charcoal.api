package context

import "fmt"
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
	Client models.Client
	User models.User
}

// ParamIntVal
//
// returns either the integer representation of a given url parameter
// or an error indicating the value was unable to be converted to the
// integer.
func (runtime *Miritos) ParamInt(name string) (int, error) {
	param := runtime.Param(name)

	if len(param) == 0 {
		return -1, fmt.Errorf("BAD_INT_VAL")
	}

	if value, err := strconv.Atoi(param); err == nil {
		return value, err
	}

	return -1, fmt.Errorf("BAD_INT_VAL")
}

// RequestHeader
// 
// helper function to return a request header from the request.
func (runtime *Miritos) RequestHeader(name string) string {
	request := runtime.Request()
	headers := request.Header()
	return headers.Get(name)
}


// Result
//
// appends a result into the runtime's result collection
func (runtime *Miritos) AddResult(result Result) {
	runtime.Results = append(runtime.Results, result)
}

// SetMeta
func (runtime *Miritos) AddMeta(key string, value interface{}) {
	runtime.Meta[key] = value
}

// ErrorOut
func (runtime *Miritos) ErrorOut(err error) error {
	runtime.Errors = append(runtime.Errors, err)
	return nil
}

// PersistFile
func (runtime *Miritos) PersistFile(target File, mime string) (models.File, error) {
	temp, err := runtime.FS.Upload(target, mime)

	if err != nil {
		return temp, err
	}

	if err := runtime.DB.Create(&temp).Error; err != nil {
		return temp, err
	}

	return temp, err
}

// Blueprint
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


// Finish
func (runtime *Miritos) Finish() error {
	runtime.Meta["time"] = time.Now()
	runtime.DB.Close()

	if runtime.Response().Committed() {
		runtime.Logger().Infof("miritos response already sent, avoiding duplicate...")
		return nil
	}

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
		Results []interface{} `json:"results"`
	}{runtime.Meta, "SUCCESS", runtime.Results.Apply()}

	return runtime.JSON(http.StatusOK, response)
}
