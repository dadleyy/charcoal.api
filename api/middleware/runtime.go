package middleware

import "github.com/kataras/iris"
import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/db"

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
func Runtime(ctx *iris.Context) {
	// attempt to connect to the mysql database
	client, err := db.Get()

	// if there was an issue opening the connection, send a 500 error
	if err != nil {
		ctx.Panic()
		return
	}

	// after the all middleware has finished, be sure to close our db connection
	defer client.Close()

	// inject our runtime into the user context for this request
	ctx.Set("runtime", api.Runtime{client})

	// move on now that it is exposed
	ctx.Next()
}
