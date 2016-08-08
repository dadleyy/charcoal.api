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
// runtime, ok := context.Get("runtime").(api.Runtime)
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

	// if there was an issue opening the connection, send a 500 error
	if err != nil {
		context.Panic()
		return
	}

	runtime = api.Runtime{Bucket: &api.Bucket{}, DB: client}

	// after the all middleware has finished, be sure to close our db connection
	defer runtime.Finish(context)

	// inject our runtime into the user context for this request
	context.Set("runtime", &runtime)

	// move on now that it is exposed
	context.Next()
}
