package testutils

import "os"
import "io"
import "net/http"
import "github.com/joho/godotenv"

import "github.com/labstack/gommon/log"
import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/activity"

type TestRouteUtil struct {
	Database *db.Connection
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

	database, _ := db.Open(dbconf)
	logger := log.New("miritos")
	queue := make(chan activity.Message)

	stub, _ := http.NewRequest(method, real, reader)

	stub.Header.Add("Content-Type", contenttype)

	server := net.ServerRuntime{logger, dbconf, queue, nil}
	route := net.Route{Method: method, Path: template}
	params, _ := route.Match(method, real)
	request, _ := server.Request(stub, &params)

	return &TestRouteUtil{database, server, request}
}
