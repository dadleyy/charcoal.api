package main

import "flag"
import "github.com/kataras/iris"
import _ "github.com/jinzhu/gorm/dialects/mysql"
import _ "github.com/joho/godotenv/autoload"

import "github.com/sizethree/meritoss.api/api/routes"
import "github.com/sizethree/meritoss.api/api/middleware"
import "github.com/sizethree/meritoss.api/api/routes/oauth"

func main() {
	flag.Parse()

	iris.UseFunc(middleware.Logger)
	iris.UseFunc(middleware.Runtime)
	iris.UseFunc(middleware.ClientAuthentication)

	iris.Get("/oauth/github", oauth.Github)

	iris.Get("/users", middleware.RequireAuth, middleware.Blueprints, routes.FindUsers)
	iris.Post("/users", routes.CreateUser)
	iris.Patch("/users/:id", middleware.RequireAuth, routes.UpdateUser)

	iris.Get("/clients", middleware.Blueprints, routes.FindClients)
	iris.Post("/clients", routes.CreateClient)

	iris.Get("/proposals", middleware.Blueprints, routes.FindProposals)
	iris.Post("/proposals", middleware.RequireAuth, routes.CreateProposal)
	// iris.Patch("/users/:id", users.Update)

	iris.Listen(":8080")
}
