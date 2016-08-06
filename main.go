package main

import "flag"
import "github.com/kataras/iris"
import _ "github.com/jinzhu/gorm/dialects/mysql"
import _ "github.com/joho/godotenv/autoload"

import "github.com/sizethree/meritoss.api/api/middleware"
import "github.com/sizethree/meritoss.api/api/routes/users"
import "github.com/sizethree/meritoss.api/api/routes/oauth"


func main() {
	flag.Parse()

	iris.UseFunc(middleware.Logger)
	iris.UseFunc(middleware.Runtime)
	iris.UseFunc(middleware.JsonAPI)

	iris.Get("/oauth/github", oauth.Github)

	iris.Get("/users", middleware.Blueprints, users.Index)
	iris.Post("/users", users.Create)
	iris.Patch("/users/:id", users.Update)

	iris.Listen(":8080")
}
