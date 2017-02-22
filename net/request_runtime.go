package net

import "net/http"

import "github.com/jinzhu/gorm"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/activity"
import "github.com/dadleyy/charcoal.api/services"
import "github.com/dadleyy/charcoal.api/filestore"

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

func (runtime *RequestRuntime) Blueprint(scopes ...*gorm.DB) Blueprint {
	cursor := runtime.DB

	if len(scopes) >= 1 {
		cursor = scopes[0]
	}

	return Blueprint{cursor, runtime.URL.Query()}
}

func (runtime *RequestRuntime) Close() {
	runtime.DB.Close()
}
