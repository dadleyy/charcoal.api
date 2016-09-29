package main

import "os"
import "flag"

import _ "github.com/joho/godotenv/autoload"
import "github.com/sizethree/miritos.api/routes"
import "github.com/sizethree/miritos.api/middleware"

func main() {
	flag.Parse()
	port := os.Getenv("PORT")

	if len(port) < 1 {
		port = "8080"
	}

	server := Server()

	server.Use(middleware.Inject)

	server.GET("/system", routes.System)
	server.GET("/auth", routes.PrintAuth, middleware.RequireUser)
	server.GET("/auth/tokens", routes.PrintClientTokens, middleware.InjectClient)

	google := server.Group("/oauth/google")

	google.GET("/prompt", routes.GoogleOauthRedirect)
	google.GET("/auth", routes.GoogleOauthReceiveCode)

	server.POST("/users", routes.CreateUser, middleware.InjectClient)
	server.GET("/users", routes.FindUser, middleware.InjectClient)
	server.PATCH("/users/:id", routes.UpdateUser, middleware.InjectClient)

	server.POST("/photos", routes.CreatePhoto, middleware.RequireUser)
	server.GET("/photos", routes.FindPhotos, middleware.RequireUser)
	server.PATCH("/photos/:id", routes.UpdatePhoto, middleware.RequireUser)


	server.Logger().Infof("starting server on port %s", port)

	Run(server, port)
}
