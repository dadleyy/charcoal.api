package testutils

import "os"
import "io"
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

func New(method, template, real, contenttype string, reader io.Reader) *TestRouteUtil {
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
	acts := make(chan activity.Message)
	socks := make(chan activity.Message)

	stub, _ := http.NewRequest(method, real, reader)

	stub.Header.Add("Content-Type", contenttype)

	server := net.ServerRuntime{logger, net.RuntimeConfig{dbconf}, acts, socks, nil}
	route := net.Route{Method: method, Path: template}
	params, _ := route.Match(method, real)
	request, _ := server.Request(stub, &params)

	return &TestRouteUtil{database, server, request}
}
