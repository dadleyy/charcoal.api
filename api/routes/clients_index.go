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
		runtime.Errors = append(runtime.Errors, err)
		context.Next()
		return
	}

	for _, client := range result {
		runtime.Results = append(runtime.Results, client)
	}

	runtime.Meta.Total = total

	context.Next()
}
