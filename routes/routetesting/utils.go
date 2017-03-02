package routetesting

import "os"
import "io"
import "bytes"
import "net/http"
import "github.com/joho/godotenv"

import "github.com/jinzhu/gorm"
import "github.com/labstack/gommon/log"
import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/dadleyy/charcoal.api/db"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/activity"

type TestRouteUtil struct {
	Database *gorm.DB
	Server   net.ServerRuntime
	Request  net.RequestRuntime
}

func NewFind(template string) *TestRouteUtil {
	reader := bytes.NewBuffer([]byte{})
	return NewRequest("GET", template, template, reader)
}

func NewPost(template string, reader io.Reader) *TestRouteUtil {
	return NewRequest("POST", template, template, reader)
}

func NewPatch(template string, real string, reader io.Reader) *TestRouteUtil {
	return NewRequest("PATCH", template, real, reader)
}

func NewRequest(method string, template string, real string, reader io.Reader) *TestRouteUtil {
	_ = godotenv.Load("../.env")

	dbconf := db.Config{
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DEBUG") == "true",
	}

	database, _ := gorm.Open("mysql", dbconf.String())
	logger := log.New("miritos")

	streams := map[string](chan activity.Message){
		"activity": make(chan activity.Message),
		"games":    make(chan activity.Message),
		"sockets":  make(chan activity.Message),
	}

	close(streams["activity"])
	close(streams["games"])
	close(streams["sockets"])

	stub, _ := http.NewRequest(method, real, reader)

	stub.Header.Add("Content-Type", "application/json")

	server := net.ServerRuntime{logger, net.RuntimeConfig{dbconf}, streams, nil}
	route := net.Route{Method: method, Path: template}
	params, _ := route.Match(method, real)
	request, _ := server.Request(stub, &params)

	return &TestRouteUtil{database, server, request}
}
