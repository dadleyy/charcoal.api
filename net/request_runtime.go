package net

import "fmt"
import "strings"
import "net/http"

import "github.com/jinzhu/gorm"
import "github.com/albrow/forms"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/defs"
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

	streams map[string](chan activity.Message)
	bucket  ResponseBucket
}

func (runtime *RequestRuntime) Form() (*forms.Data, error) {
	body, err := forms.Parse(runtime.Request)
	return body, err
}

func (runtime *RequestRuntime) IsAdmin() bool {
	uman := services.UserManager{runtime.DB, runtime.Logger}
	return uman.IsAdmin(&runtime.User)
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

func (runtime *RequestRuntime) AddError(list ...error) error {
	if len(list) == 0 {
		return nil
	}

	msgs := make([]string, 0, len(list))

	for _, e := range list {
		runtime.bucket.errors = append(runtime.bucket.errors, e)
		msgs = append(msgs, e.Error())
	}

	return fmt.Errorf(strings.Join(msgs, " | "))
}

func (runtime *RequestRuntime) appendErrors(identifier string, values ...string) error {
	result := make([]error, 0, len(values))

	for _, f := range values {
		result = append(result, fmt.Errorf("%s:%s", identifier, f))
	}

	return runtime.AddError(result...)
}

func (runtime *RequestRuntime) LogicError(reasons ...string) error {
	return runtime.appendErrors("reason", reasons...)
}

func (runtime *RequestRuntime) FieldError(fields ...string) error {
	return runtime.appendErrors("field", fields...)
}

func (runtime *RequestRuntime) ServerError() error {
	return runtime.AddError(fmt.Errorf("reason:server-error"))
}

func (runtime *RequestRuntime) SetMeta(key string, val interface{}) {
	runtime.bucket.meta[key] = val
}

func (runtime *RequestRuntime) SetTotal(total int) {
	runtime.SetMeta("count", total)
}

func (runtime *RequestRuntime) Publish(msg activity.Message) {
	identifiers := strings.Split(msg.Verb, ":")

	if len(identifiers) != 2 {
		runtime.Debugf("[runtime] invalid message identifier: %s", msg.Verb)
		return
	}

	stream, ok := runtime.streams[identifiers[0]]

	if ok != true || stream == nil {
		runtime.Warnf("[runtime] invalid message identifier: %s", msg.Verb)
		return
	}

	stream <- msg
}

func (runtime *RequestRuntime) Photos() services.PhotoSaver {
	return services.PhotoSaver{runtime.DB, runtime.FileSaver}
}

func (runtime *RequestRuntime) Game(id uint) (*services.GameManager, error) {
	g := models.Game{}

	if e := runtime.First(&g, id).Error; e != nil {
		return nil, e
	}

	streams := map[string](chan<- activity.Message){
		defs.SocketsStreamIdentifier: runtime.streams[defs.SocketsStreamIdentifier],
		defs.GamesStreamIdentifier:   runtime.streams[defs.GamesStreamIdentifier],
	}

	manager := services.GameManager{runtime.DB, runtime.Logger, streams, g}

	return &manager, nil
}

func (runtime *RequestRuntime) Blueprint(scopes ...*gorm.DB) Blueprint {
	cursor := runtime.DB

	if len(scopes) >= 1 {
		cursor = scopes[0]
	}

	return Blueprint{cursor, runtime.Logger, runtime.URL.Query()}
}

func (runtime *RequestRuntime) Close() {
	runtime.DB.Close()
}
