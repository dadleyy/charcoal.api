package main

import "flag"
import "github.com/kataras/iris"
import _ "github.com/jinzhu/gorm/dialects/mysql"
import _ "github.com/joho/godotenv/autoload"

import "github.com/sizethree/meritoss.api/routes"
import "github.com/sizethree/meritoss.api/middleware"
import "github.com/sizethree/meritoss.api/routes/oauth"

func main() {
	flag.Parse()

	iris.UseFunc(middleware.Logger)
	iris.UseFunc(middleware.InjectRuntime)
	iris.UseFunc(middleware.ClientAuthentication)

	iris.Get("/oauth/github", oauth.Github)

	iris.Get("/users", middleware.RequireAuth, middleware.InjectBlueprint, routes.FindUsers)
	iris.Post("/users", routes.CreateUser)
	iris.Patch("/users/:id", middleware.RequireAuth, routes.UpdateUser)

	iris.Get("/clients", middleware.InjectBlueprint, routes.FindClients)
	iris.Post("/clients", middleware.RequireAuth, routes.CreateClient)

	iris.Get("/proposals", middleware.InjectBlueprint, routes.FindProposals)
	iris.Post("/proposals", middleware.RequireAuth, routes.CreateProposal)
	// iris.Patch("/users/:id", users.Update)

	iris.Get("/positions", middleware.InjectBlueprint, routes.FindPositions)
	iris.Post("/positions", middleware.RequireAuth, middleware.InjectBlueprint, routes.CreatePosition)

	iris.Get("/clienttokens", middleware.RequireAuth, routes.FindClientTokens)
	iris.Post("/clienttokens", middleware.RequireAuth, routes.CreateClientToken)

	iris.Listen(":8080")
}
