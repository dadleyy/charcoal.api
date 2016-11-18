package main

import "os"
import "fmt"

import "github.com/joho/godotenv"
import "github.com/labstack/gommon/log"

import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/routes"
import "github.com/sizethree/miritos.api/activity"
import "github.com/sizethree/miritos.api/middleware"

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("bad env: %s\n", err.Error())
		return
	}

	dbconf := db.Config{
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DEBUG") == "true",
	}

	port := os.Getenv("PORT")

	if len(port) < 1 {
		port = "8080"
	}

	database, err := db.Open(dbconf)

	if err != nil {
		panic(err)
	}

	// create the logger that will be shared by the server and the activity processor
	logger := log.New("miritos")
	logger.SetLevel(0)
	logger.SetHeader("[${time_rfc3339} ${level} ${short_file}]")

	// create the channel that will be used by the server runtime and activity processor
	stream := make(chan activity.Message, 100)

	// create our multiplexer and add our routes
	mux := net.Multiplexer{}

	mux.Use(middleware.InjectClient)
	mux.Use(middleware.InjectUser)

	mux.GET("/system", routes.PrintSystem)
	mux.PATCH("/system", routes.UpdateSystem, middleware.RequireUser, middleware.RequireAdmin)

	mux.GET("/system/domains", routes.FindSystemEmailDomains, middleware.RequireUser, middleware.RequireAdmin)
	mux.POST("/system/domains", routes.CreateSystemEmailDomain, middleware.RequireUser, middleware.RequireAdmin)
	mux.DELETE("/system/domains/:id", routes.DestroySystemEmailDomain, middleware.RequireUser, middleware.RequireAdmin)

	mux.GET("/auth/user", routes.PrintAuth, middleware.RequireUser)
	mux.GET("/auth/roles", routes.PrintUserRoles, middleware.RequireUser)
	mux.GET("/auth/tokens", routes.PrintClientTokens, middleware.RequireClient)

	mux.GET("/activity", routes.FindActivity)

	// special route - returns all active activities based on their display schedules.
	mux.GET("/activity/live", routes.FindLiveActivity, middleware.RequireClient)

	mux.POST("/callbacks/mailgun", routes.MailgunUploadHook)

	mux.GET("/oauth/google/prompt", routes.GoogleOauthRedirect)
	mux.GET("/oauth/google/auth", routes.GoogleOauthReceiveCode)

	mux.GET("/user-roles", routes.FindRoles, middleware.RequireClient)

	// clients
	mux.GET("/clients", routes.FindClients, middleware.RequireClient)
	mux.POST("/clients", routes.CreateClient, middleware.RequireUser)

	mux.GET("/google-accounts", routes.FindGoogleAccounts, middleware.RequireUser)

	// client admins
	mux.GET("/client-admins", routes.FindClientAdmins, middleware.RequireClient, middleware.RequireUser)
	mux.POST("/client-admins", routes.CreateClientAdmin, middleware.RequireClient, middleware.RequireUser)

	// client tokens
	mux.GET("/client-tokens", routes.FindClientTokens, middleware.RequireUser, middleware.RequireAdmin)
	mux.POST("/client-tokens", routes.CreateClientToken, middleware.RequireClient, middleware.RequireUser)

	mux.GET("/display-schedules", routes.FindDisplaySchedules, middleware.RequireClient)
	mux.PATCH("/display-schedules/:id", routes.UpdateDisplaySchedule, middleware.RequireUser)

	mux.GET("/user-role-mappings", routes.FindUserRoleMappings, middleware.RequireUser)
	mux.POST("/user-role-mappings", routes.CreateUserRoleMapping, middleware.RequireUser, middleware.RequireAdmin)
	mux.DELETE("/user-role-mappings/:id", routes.DestroyUserRoleMapping, middleware.RequireUser, middleware.RequireAdmin)

	mux.POST("/instagram", routes.CreateInstagramPost, middleware.RequireClient)
	mux.GET("/instagram", routes.FindInstagramPosts, middleware.RequireClient)

	mux.GET("/users", routes.FindUser, middleware.RequireClient)
	mux.POST("/users", routes.CreateUser, middleware.RequireClient)
	mux.PATCH("/users/:id", routes.UpdateUser, middleware.RequireUser)

	mux.POST("/photos", routes.CreatePhoto, middleware.RequireClient)
	mux.GET("/photos", routes.FindPhotos, middleware.RequireClient)
	mux.GET("/photos/:id/view", routes.ViewPhoto, middleware.RequireClient)
	mux.DELETE("/photos/:id", routes.DestroyPhoto, middleware.RequireUser)

	// create the server runtime and the activity processor runtime
	runtime := net.ServerRuntime{logger, dbconf, stream, &mux}
	processor := activity.Processor{logger, database, stream}
	server := net.Server{nil, &runtime}

	// start the server & processor
	server.Logger().Debugf(fmt.Sprintf("starting"))
	go processor.Begin()
	server.Run(fmt.Sprintf(":%s", port))
}
