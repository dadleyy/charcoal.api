package routetesting

import "io"
import "bytes"
import "strconv"
import "net/url"
import "net/http"

import "github.com/jinzhu/gorm"
import "github.com/labstack/gommon/log"
import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/activity"
import "github.com/dadleyy/charcoal.api/testutils"

type TestRouteUtil struct {
	Database *gorm.DB
	Server   net.ServerRuntime
	Request  *net.RequestRuntime
	Streams  map[string](chan activity.Message)
}

type TestRouteParams struct {
	url.Values
}

func (params *TestRouteParams) IntParam(key string) (int, bool) {
	v := params.Get(key)
	i, e := strconv.Atoi(v)
	return i, e == nil
}

func (params *TestRouteParams) StringParam(key string) (string, bool) {
	return params.Get(key), false
}

func NewFind(params *TestRouteParams) *TestRouteUtil {
	reader := bytes.NewBuffer([]byte{})
	return NewRequest("GET", params, reader)
}

func NewPost(params *TestRouteParams, reader io.Reader) *TestRouteUtil {
	return NewRequest("POST", params, reader)
}

func NewPatch(params *TestRouteParams, reader io.Reader) *TestRouteUtil {
	return NewRequest("PATCH", params, reader)
}

func NewRequest(method string, params *TestRouteParams, reader io.Reader) *TestRouteUtil {
	database := testutils.NewDB()

	logger := log.New("miritos")

	streams := map[string](chan activity.Message){
		"activity": make(chan activity.Message),
		"games":    make(chan activity.Message),
		"sockets":  make(chan activity.Message),
	}

	stub, _ := http.NewRequest(method, "/blah", reader)

	stub.Header.Add("Content-Type", "application/json")

	server := net.ServerRuntime{
		Logger:  logger,
		DB:      database,
		Streams: streams,
		Mux:     nil,
	}

	request := server.Request(stub, params)

	return &TestRouteUtil{database, server, request, streams}
}
