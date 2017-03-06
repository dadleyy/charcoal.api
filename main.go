package main

import "os"
import "fmt"
import "net/http"

import "github.com/joho/godotenv"
import "github.com/labstack/gommon/log"

import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/dadleyy/charcoal.api/db"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/routes"
import "github.com/dadleyy/charcoal.api/activity"
import "github.com/dadleyy/charcoal.api/middleware"

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

	if err != nil {
		panic(err)
	}

	// create the logger that will be shared by the server and the activity processor
	logger := log.New("miritos")
	logger.SetLevel(0)
	logger.SetHeader("[${time_rfc3339} ${level} ${short_file}]")

	streams := map[string](chan activity.Message){
		"activity": make(chan activity.Message, 100),
		"sockets":  make(chan activity.Message, 100),
		"games":    make(chan activity.Message, 100),
	}

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
	mux.POST("/auth", routes.PasswordLogin, middleware.RequireClient)

	mux.GET("/activity", routes.FindActivity)

	mux.GET("/activity/live", routes.FindLiveActivity, middleware.RequireClient)

	mux.POST("/callbacks/mailgun", routes.MailgunUploadHook)

	mux.GET("/oauth/google/prompt", routes.GoogleOauthRedirect)
	mux.GET("/oauth/google/auth", routes.GoogleOauthReceiveCode)

	mux.GET("/user-roles", routes.FindRoles, middleware.RequireClient)

	mux.GET("/clients", routes.FindClients, middleware.RequireClient)
	mux.POST("/clients", routes.CreateClient, middleware.RequireUser)
	mux.PATCH("/clients/:id", routes.UpdateClient, middleware.RequireUser)

	mux.GET("/client-admins", routes.FindClientAdmins, middleware.RequireClient, middleware.RequireUser)
	mux.POST("/client-admins", routes.CreateClientAdmin, middleware.RequireClient, middleware.RequireUser)
	mux.DELETE("/client-admins/:id", routes.DeleteClientAdmin, middleware.RequireClient, middleware.RequireUser)

	mux.GET("/client-tokens", routes.FindClientTokens, middleware.RequireUser, middleware.RequireAdmin)
	mux.POST("/client-tokens", routes.CreateClientToken, middleware.RequireClient, middleware.RequireUser)

	mux.GET("/google-accounts", routes.FindGoogleAccounts, middleware.RequireUser)

	mux.GET("/display-schedules", routes.FindDisplaySchedules, middleware.RequireClient)
	mux.PATCH("/display-schedules/:id", routes.UpdateDisplaySchedule, middleware.RequireUser)

	mux.GET("/user-role-mappings", routes.FindUserRoleMappings, middleware.RequireUser)
	mux.POST("/user-role-mappings", routes.CreateUserRoleMapping, middleware.RequireUser, middleware.RequireAdmin)
	mux.DELETE("/user-role-mappings/:id", routes.DestroyUserRoleMapping, middleware.RequireUser, middleware.RequireAdmin)

	mux.GET("/users", routes.FindUser, middleware.RequireClient)
	mux.POST("/users", routes.CreateUser, middleware.RequireClient)
	mux.PATCH("/users/:id", routes.UpdateUser, middleware.RequireUser)

	mux.POST("/games", routes.CreateGame, middleware.RequireUser)
	mux.GET("/games", routes.FindGames, middleware.RequireUser)
	mux.DELETE("/games/:id", routes.DestroyGame, middleware.RequireUser)
	mux.PATCH("/games/:id", routes.UpdateGame, middleware.RequireUser)

	mux.POST("/game-rounds", routes.CreateGameRound, middleware.RequireUser)
	mux.GET("/game-rounds", routes.FindGameRounds, middleware.RequireUser)
	mux.PATCH("/game-rounds/:id", routes.UpdateGameRound, middleware.RequireUser)
	mux.DELETE("/game-rounds/:id", routes.DestroyGameRound, middleware.RequireUser)

	mux.POST("/game-memberships", routes.CreateGameMembership, middleware.RequireUser)
	mux.GET("/game-memberships", routes.FindGameMemberships, middleware.RequireUser)
	mux.DELETE("/game-memberships/:id", routes.DestroyGameMembership, middleware.RequireUser)

	mux.POST("/photos", routes.CreatePhoto, middleware.RequireClient)
	mux.GET("/photos", routes.FindPhotos, middleware.RequireClient)
	mux.GET("/photos/:id/view", routes.ViewPhoto, middleware.RequireClient)
	mux.DELETE("/photos/:id", routes.DestroyPhoto, middleware.RequireUser)

	// create the server runtime and the activity processor runtime
	runtime := net.ServerRuntime{
		Logger:  logger,
		Config:  net.RuntimeConfig{dbconf},
		Mux:     &mux,
		Streams: streams,
	}

	processors := []activity.BackgroundProcessor{
		&activity.ActivityProcessor{Logger: logger, Stream: streams["activity"]},
		&activity.GameProcessor{Logger: logger, Stream: streams["games"], Exhaust: nil},
	}

	if os.Getenv("SOCKETS_ENABLED") == "true" {
		websock := net.SocketRuntime{logger, streams["sockets"]}
		http.Handle("/socket/", &websock)
	}

	http.Handle("/", &runtime)

	for _, processor := range processors {
		go processor.Begin(activity.ProcessorConfig{dbconf})
	}

	logger.Debugf(fmt.Sprintf("starting on port: %s", port))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
