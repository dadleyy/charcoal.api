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
}

func (runtime *RequestRuntime) Form() (*forms.Data, error) {
	body, err := forms.Parse(runtime.Request)
	return body, err
}

func (runtime *RequestRuntime) IsAdmin() bool {
	uman := services.UserManager{runtime.DB, runtime.Logger}
	return uman.IsAdmin(&runtime.User)
}

func (runtime *RequestRuntime) Redirect(url string) *ResponseBucket {
	return &ResponseBucket{Redirect: url}
}

func (runtime *RequestRuntime) Proxy(url string) *ResponseBucket {
	return &ResponseBucket{Proxy: url}
}

func (runtime *RequestRuntime) SendErrors(list ...error) *ResponseBucket {
	return &ResponseBucket{Errors: list}
}

func (runtime *RequestRuntime) appendErrors(identifier string, values ...string) *ResponseBucket {
	result := make([]error, 0, len(values))

	for _, f := range values {
		result = append(result, fmt.Errorf("%s:%s", identifier, f))
	}

	return runtime.SendErrors(result...)
}

func (runtime *RequestRuntime) LogicError(reasons ...string) *ResponseBucket {
	return runtime.appendErrors("reason", reasons...)
}

func (runtime *RequestRuntime) FieldError(fields ...string) *ResponseBucket {
	return runtime.appendErrors("field", fields...)
}

func (runtime *RequestRuntime) ServerError() *ResponseBucket {
	return runtime.SendErrors(fmt.Errorf("reason:server-error"))
}

func (runtime *RequestRuntime) SendResults(total int, results interface{}) *ResponseBucket {
	meta := map[string]interface{}{"total": total}
	return &ResponseBucket{Results: results, Meta: meta}
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
}
