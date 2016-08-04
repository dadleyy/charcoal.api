package main

import "flag"
import "github.com/kataras/iris"
import _ "github.com/jinzhu/gorm/dialects/mysql"
import _ "github.com/joho/godotenv/autoload"

import "github.com/meritoss/meritoss.api/api/middleware"
import "github.com/meritoss/meritoss.api/api/routes/users"


func main() {
	flag.Parse()

	iris.UseFunc(middleware.Logger)
	iris.UseFunc(middleware.Runtime)
	iris.Get("/users", middleware.Blueprints, users.Index)
	iris.Post("/users", users.Create)
	iris.Listen(":8080")
}
