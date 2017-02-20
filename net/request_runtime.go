package net

import "strconv"
import "strings"
import "net/http"

import "github.com/jinzhu/gorm"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/activity"
import "github.com/dadleyy/charcoal.api/services"
import "github.com/dadleyy/charcoal.api/filestore"

const DEFAULT_BLUEPRINT_LIMIT = 100

type RequestRuntime struct {
	*http.Request
	*UrlParams
	*log.Logger
	*gorm.DB
	filestore.FileSaver

	Client models.Client
	User   models.User

	queue  chan activity.Message
	bucket ResponseBucket
}

func (runtime *RequestRuntime) IsAdmin() bool {
	uman := services.UserManager{runtime.DB}
	return uman.IsAdmin(&runtime.User) && runtime.Client.System == true
}

func (runtime *RequestRuntime) AddResult(r Result) {
	runtime.bucket.results = append(runtime.bucket.results, r)
}

func (runtime *RequestRuntime) Redirect(url string) {
	runtime.bucket.redirect = url
}

func (runtime *RequestRuntime) Proxy(url string) {
	runtime.bucket.proxy = url
}

func (runtime *RequestRuntime) AddError(e error) error {
	runtime.bucket.errors = append(runtime.bucket.errors, e)
	return e
}

func (runtime *RequestRuntime) SetMeta(key string, val interface{}) {
	runtime.bucket.meta[key] = val
}

func (runtime *RequestRuntime) Publish(msg activity.Message) {
	runtime.queue <- msg
}

func (runtime *RequestRuntime) Photos() services.PhotoSaver {
	return services.PhotoSaver{runtime.DB, runtime.FileSaver}
}

func (runtime *RequestRuntime) Blueprint() Blueprint {
	result := Blueprint{limit: DEFAULT_BLUEPRINT_LIMIT, page: 0}

	values := runtime.URL.Query()

	if page, ok := values["page"]; ok && len(page) == 1 {
		ipage, err := strconv.Atoi(page[0])

		if err == nil {
			result.page = ipage
		}
	}

	if limit, ok := values["limit"]; ok && len(limit) == 1 {
		ilimit, err := strconv.Atoi(limit[0])

		if err == nil {
			result.limit = ilimit
		}
	}

	for key, values := range values {
		filterable := strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]")

		if filterable == false || len(values) != 1 {
			continue
		}

		value := values[0]

		if err := result.Filter(key, value); err != nil {
			runtime.AddError(err)
		}
	}

	return result
}

func (runtime *RequestRuntime) Close() {
	runtime.DB.Close()
}

func (runtime *RequestRuntime) Cursor(start interface{}) *gorm.DB {
	return runtime.Model(start)
}

func (runtime *RequestRuntime) Database() *gorm.DB {
	return runtime.DB
}
