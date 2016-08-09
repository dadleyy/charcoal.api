package middleware

import "github.com/kataras/iris"
import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/db"

// Runtime
//
// Definess a middleware function that will inject an instance of `api.Runtime` into 
// each request so that the handler can access things like the orm by retreiving the 
// value from the context: 
//
// runtime, ok := context.Get("runtime").(*api.Runtime)
// 
// where api.Runtime is a struct that could look like:
//
// type Runtime struct {
//   DB *gorm.DB
// }
// 
// In this example, we're using the "github.com/jinzhu/gorm" package to handle db
// related communication and modeling.
func Runtime(context *iris.Context) {
	var runtime api.Runtime

	// attempt to connect to the mysql database
	client, err := db.Get()

	// if there was an issue opening the connection, send a 500 error - nothing we 
	// can do here to solve that problem
	if err != nil {
		context.Panic()
		return
	}

	// initialize the runtime struct with the gorm client; all other properties will
	// be default initialized
	runtime = api.Runtime{DB: client}

	// after the all middleware has completed, let the runtime do whatever it needs to
	// do in order to complete this request.
	defer runtime.Render(context)

	// inject our runtime into the user context for this request
	context.Set("runtime", &runtime)

	// move on now that it is exposed
	context.Next()
}
