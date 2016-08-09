package routes

import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/dal"

func FindUsers(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*api.Runtime)

	if !ok {
		glog.Error("bad runtime")
		context.Panic()
		context.StopExecution()
		return
	}

	blueprint, _ := context.Get("blueprint").(*api.Blueprint)

	result, total, err := dal.FindUser(&runtime.DB, blueprint)

	if err != nil {
		runtime.Error(err)
		return
	}

	for _, user := range result {
		runtime.Result(user)
	}

	runtime.Meta("total", total)
	context.Next()
}
