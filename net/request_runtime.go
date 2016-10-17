package net

import "strconv"
import "strings"
import "net/http"

import "github.com/labstack/gommon/log"
import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/activity"
import "github.com/sizethree/miritos.api/filestore"

const DEFAULT_BLUEPRINT_LIMIT = 100

type RequestRuntime struct {
	*http.Request
	*UrlParams
	database *db.Connection
	log *log.Logger
	queue chan activity.Message
	Client models.Client
	User models.User
	bucket ResponseBucket
	store filestore.FileSaver
}

func (runtime *RequestRuntime) Errorf(format string, args ...interface{}) {
	runtime.log.Errorf(format, args...)
}

func (runtime *RequestRuntime) Warnf(format string, args ...interface{}) {
	runtime.log.Warnf(format, args...)
}

func (runtime *RequestRuntime) Infof(format string, args ...interface{}) {
	runtime.log.Infof(format, args...)
}

func (runtime *RequestRuntime) Debugf(format string, args ...interface{}) {
	runtime.log.Debugf(format, args...)
}

func (runtime *RequestRuntime) Printf(format string, args ...interface{}) {
	runtime.log.Printf(format, args...)
}

func (runtime *RequestRuntime) AddResult(r Result) {
	runtime.bucket.results = append(runtime.bucket.results, r)
}

func (runtime *RequestRuntime) PersistFile(file filestore.File, mime string) (models.File, error) {
	ofile, err := runtime.store.Upload(file, mime)

	if err != nil {
		return models.File{}, err
	}

	if err := runtime.Database().Create(&ofile).Error; err != nil {
		return models.File{}, err
	}

	return ofile, nil
}

func (runtime *RequestRuntime) AddError(e error) error {
	runtime.bucket.errors = append(runtime.bucket.errors, e)
	return nil
}

func (runtime *RequestRuntime) SetMeta(key string, val interface{}) {
	runtime.bucket.meta[key] = val
}

func (runtime *RequestRuntime) Publish(msg activity.Message) {
	runtime.queue <- msg
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

func (runtime *RequestRuntime) Database() *db.Connection {
	return runtime.database
}

