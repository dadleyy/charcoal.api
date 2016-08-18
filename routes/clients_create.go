package routes

import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/miritos.api/dal"
import "github.com/sizethree/miritos.api/middleware"


// CreateClient
//
// request handler for POST /clients
// 
// will attempt to load in json data as a `ClientFacade` and use the `dal.CreateClient`
// function to persist that information to the clients table
func CreateClient(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*middleware.Runtime)

	if !ok {
		glog.Error("bad runtime")
		context.Panic()
		context.StopExecution()
		return
	}

	// prepare our facade for iris to load into
	var target dal.ClientFacade

	// if we fail at reading json into the structure prescribed by the data access layer,
	// add the error we receive to our runtime and continue on.
	if e := context.ReadJSON(&target); e != nil {
		runtime.Error(e)
		context.Next()
		return
	}

	// attempt to create the new client
	client, e := dal.CreateClient(&runtime.DB, &target);

	// if we fail at creating a client - move on
	if e != nil {
		runtime.Error(e)
		context.Next()
		return
	}

	runtime.Result(client)
	runtime.Meta("total", 1)
	glog.Infof("created client %d\n", client.ID)
	context.Next()
}
