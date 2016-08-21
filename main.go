package main

import "os"
import "fmt"
import "flag"
import "github.com/dadleyy/iris"
import _ "github.com/jinzhu/gorm/dialects/mysql"
import _ "github.com/joho/godotenv/autoload"

import "github.com/sizethree/miritos.api/routes"
import "github.com/sizethree/miritos.api/middleware"
import "github.com/sizethree/miritos.api/routes/oauth"

func main() {
	flag.Parse()
	port := os.Getenv("PORT")

	if len(port) < 1 {
		port = "8080"
	}

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
	iris.Patch("/proposals/:id", middleware.RequireAuth, routes.UpdateProposal)

	iris.Get("/positions", middleware.InjectBlueprint, routes.FindPositions)
	iris.Post("/positions", middleware.RequireAuth, routes.CreatePosition)
	iris.Patch("/positions/:id", middleware.RequireAuth, routes.UpdatePosition)

	iris.Get("/clienttokens", middleware.RequireAuth, routes.FindClientTokens)
	iris.Post("/clienttokens", middleware.RequireAuth, routes.CreateClientToken)

	iris.Listen(fmt.Sprintf(":%s", port))
}
