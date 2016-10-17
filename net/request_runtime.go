package net

import "net/http"
import "github.com/labstack/gommon/log"
import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/activity"

type RequestRuntime struct {
	*http.Request
	*UrlParams
	DatabaseConnection *db.Connection
	Log *log.Logger
	Queue chan activity.Message
	Client models.Client
	User models.User
	bucket ResponseBucket
}

func (runtime *RequestRuntime) Errorf(format string, args ...interface{}) {
	runtime.Log.Errorf(format, args...)
}

func (runtime *RequestRuntime) Warnf(format string, args ...interface{}) {
	runtime.Log.Warnf(format, args...)
}

func (runtime *RequestRuntime) Infof(format string, args ...interface{}) {
	runtime.Log.Infof(format, args...)
}

func (runtime *RequestRuntime) Debugf(format string, args ...interface{}) {
	runtime.Log.Debugf(format, args...)
}

func (runtime *RequestRuntime) Printf(format string, args ...interface{}) {
	runtime.Log.Printf(format, args...)
}

func (runtime *RequestRuntime) Result(r Result) {
	runtime.bucket.results = append(runtime.bucket.results, r)
}

func (runtime *RequestRuntime) Error(e error) error {
	runtime.bucket.errors = append(runtime.bucket.errors, e)
	return nil
}

func (runtime *RequestRuntime) Database() *db.Connection {
	return runtime.DatabaseConnection
}
