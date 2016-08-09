package routes

import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/dal"

func FindClients(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*api.Runtime)

	if !ok {
		glog.Error("bad runtime")
		context.Panic()
		context.StopExecution()
		return
	}

	blueprint, _ := context.Get("blueprint").(*api.Blueprint)

	result, total, err := dal.FindClients(&runtime.DB, blueprint)

	if err != nil {
		glog.Errorf("failed finding valid clients: %s\n", err.Error())
		runtime.Error(err)
		context.Next()
		return
	}

	// add each client to our results
	for _, client := range result {
		runtime.Result(client)
	}

	// add the total we received
	runtime.Meta("total", total)

	// move on
	context.Next()
}
