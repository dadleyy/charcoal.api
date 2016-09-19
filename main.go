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

	server.POST("/users", routes.CreateUser)
	server.GET("/users", routes.FindUser)
	server.PATCH("/users/:id", routes.UpdateUser)

	server.Logger().Infof("starting server on port %s", port)

	Run(server, port)
}
