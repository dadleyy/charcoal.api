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
	server.GET("/auth", routes.PrintAuth, middleware.UserAuthentication)

	google := server.Group("/oauth/google")

	google.GET("/prompt", routes.GoogleOauthRedirect)
	google.GET("/auth", routes.GoogleOauthReceiveCode)

	server.POST("/users", routes.CreateUser, middleware.ClientAuthentication)
	server.GET("/users", routes.FindUser, middleware.ClientAuthentication)
	server.PATCH("/users/:id", routes.UpdateUser, middleware.ClientAuthentication)

	server.POST("/photos", routes.CreatePhoto, middleware.ClientAuthentication)
	server.GET("/photos", routes.FindPhotos, middleware.ClientAuthentication)
	server.PATCH("/photos/:id", routes.UpdatePhoto, middleware.ClientAuthentication)


	server.Logger().Infof("starting server on port %s", port)

	Run(server, port)
}
