package routes

import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/dal"

func FindProposals(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*api.Runtime)

	if !ok {
		glog.Error("bad runtime")
		context.Panic()
		context.StopExecution()
		return
	}

	blueprint, _ := context.Get("blueprint").(*api.Blueprint)

	result, total, err := dal.FindProposals(&runtime.DB, blueprint)

	if err != nil {
		glog.Errorf("failed finding valid proposals: %s\n", err.Error())
		runtime.Errors = append(runtime.Errors, err)
		return
	}

	for _, proposal := range result {
		runtime.Results = append(runtime.Results, proposal)
	}

	runtime.Meta.Total = total

	context.Next()
}
