package main

import "os"
import "fmt"
import "flag"
import "github.com/jinzhu/gorm"
import "github.com/kataras/iris"
import _ "github.com/jinzhu/gorm/dialects/mysql"
import _ "github.com/joho/godotenv/autoload"

import "github.com/meritoss/meritoss.api/api"
import "github.com/meritoss/meritoss.api/api/routes/users"

const DSN_STR = "%v:%v@tcp(%v:%v)/%v"

// runtime
//
// this is a middleware function that will inject an instance of `api.AppContext` into 
// each request so that the handler can access things like the orm by retreiving the 
// value from the context: 
//
// runtime, ok := ctx.Get("runtime").(api.Runtime)
// 
// where api.Runtime is a struct that could look like:
//
// type Runtime struct {
//   DB *gorm.DB
// }
// 
// In this example, we're using the "github.com/jinzhu/gorm" package to handle db
// related communication and modeling.
func runtime(ctx *iris.Context) {
	// get configuration information from the environment
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	hostname := os.Getenv("DB_HOSTNAME")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	// build our data source url
	dsn := fmt.Sprintf(DSN_STR, username, password, hostname, port, database)

	// attempt to connect to the mysql database
	db, err := gorm.Open("mysql", dsn)

	// if there was an issue opening the connection, send a 500 error
	if err != nil {
		ctx.Panic()
		return
	}

	// after the all middleware has finished, be sure to close our db connection
	defer db.Close()

	// turn off gorm logging
	db.LogMode(false)

	// inject our runtime into the user context for this request
	ctx.Set("runtime", api.Runtime{db})

	// move on now that it is exposed
	ctx.Next()
}

func main() {
	flag.Parse()
	iris.UseFunc(runtime)
	iris.Get("/users", users.Index)
	iris.Listen(":8080")
}
