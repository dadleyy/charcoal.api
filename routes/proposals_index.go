package routes

import "github.com/golang/glog"
import "github.com/dadleyy/iris"

import "github.com/sizethree/miritos.api/dal"
import "github.com/sizethree/miritos.api/middleware"

func FindProposals(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*middleware.Runtime)

	if !ok {
		glog.Error("bad runtime")
		context.Panic()
		context.StopExecution()
		return
	}

	blueprint, _ := context.Get("blueprint").(*middleware.Blueprint)

	result, total, err := dal.FindProposals(&runtime.DB, blueprint)

	if err != nil {
		glog.Errorf("failed finding valid proposals: %s\n", err.Error())
		runtime.Error(err)
		return
	}

	for _, proposal := range result {
		runtime.Result(proposal)
	}

	runtime.Meta("total", total)
	context.Next()
}
