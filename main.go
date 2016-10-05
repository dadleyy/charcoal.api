package main

import "os"
import "flag"

import _ "github.com/joho/godotenv/autoload"
import "github.com/sizethree/miritos.api/routes"
import "github.com/sizethree/miritos.api/server"
import "github.com/sizethree/miritos.api/activity"
import "github.com/sizethree/miritos.api/middleware"

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

	stream := make(chan activity.Message, 100)
	dbconf := server.DatabaseConfig{dbusername, dbpassword, dbhostname, dbdatabase, dbport, dbdebug}

	app := server.NewApp()
	processor := activity.Processor{stream, dbconf}

	go processor.Begin()

	app.Use(middleware.Inject(stream, dbconf))

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
	app.GET("/photos", routes.FindPhotos, middleware.RequireUser)
	app.GET("/photos/:id/view", routes.ViewPhoto, middleware.RequireClient)

	app.GET("/activity", routes.FindActivity, middleware.RequireClient)

	app.Logger().Infof("starting app on port %s", port)

	server.RunApp(app, port)
}
