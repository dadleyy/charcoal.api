package main

import "os"
import "fmt"
import "flag"

import "github.com/labstack/echo"
import "github.com/labstack/gommon/log"
import "github.com/labstack/echo/engine/standard"

import _ "github.com/joho/godotenv/autoload"
import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/routes"
import "github.com/sizethree/miritos.api/server"
import "github.com/sizethree/miritos.api/activity"
import "github.com/sizethree/miritos.api/middleware"

const logfmt = "[${level} ${prefix} ${short_file}:${line}]"

func main() {
	flag.Parse()
	port := os.Getenv("PORT")

	if len(port) < 1 {
		port = "8080"
	}

	dbusername := os.Getenv("DB_USERNAME")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbhostname := os.Getenv("DB_HOSTNAME")
	dbport := os.Getenv("DB_PORT")
	dbdatabase := os.Getenv("DB_DATABASE")
	dbdebug := os.Getenv("DB_DEBUG") == "true"

	// create the stream that will handle our activity messages
	stream := make(chan activity.Message, 100)

	// prepare the database configuration that will be used both in the runtime that handles responses,
	// as well as the connection used by the activity processor to handle them.
	dbconf := db.Config{dbusername, dbpassword, dbhostname, dbdatabase, dbport, dbdebug}

	// prepare the logger shared by the server runtime and the activity processor
	logger := log.New("miritos")
	logger.SetHeader(logfmt)

	// create the main application and the processor that will be handling the activity messages
	app := server.App{Echo: echo.New(), Queue: stream, DBConfig: dbconf}

	processor := activity.Processor{stream, logger, dbconf}

	app.SetLogger(logger)
	app.SetLogLevel(0)

	// start the processor's goroutine
	go processor.Begin()

	// add the runtime injection middleware
	app.Use(app.Inject)

	app.GET("/system", routes.System)

	app.GET("/auth", routes.PrintAuth, middleware.RequireUser)
	app.GET("/auth/tokens", routes.PrintClientTokens, middleware.InjectClient)

	google := app.Group("/oauth/google")

	google.GET("/prompt", routes.GoogleOauthRedirect)
	google.GET("/auth", routes.GoogleOauthReceiveCode)

	app.POST("/users", routes.CreateUser, middleware.RequireClient)
	app.GET("/users", routes.FindUser, middleware.RequireClient)
	app.PATCH("/users/:id", routes.UpdateUser, middleware.RequireUser)

	app.POST("/photos", routes.CreatePhoto, middleware.RequireUser)
	app.GET("/photos", routes.FindPhotos, middleware.RequireClient)
	app.GET("/photos/:id/view", routes.ViewPhoto, middleware.RequireClient)

	app.GET("/activity", routes.FindActivity, middleware.RequireClient)

	app.GET("/displayschedules", routes.FindDisplaySchedules, middleware.RequireClient)

	app.Logger().Infof("starting app on port %s", port)

	app.Run(standard.New(fmt.Sprintf(":%s", port)))
}
