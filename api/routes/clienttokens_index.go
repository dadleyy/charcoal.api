package routes

import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/dal"

func FindClientTokens(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*api.Runtime)

	if !ok {
		glog.Error("bad runtime")
		context.Panic()
		context.StopExecution()
		return
	}

	user := runtime.User.ID
	glog.Infof("looking up client tokens for user %d\n")

	results, total, e := dal.ClientTokensForUser(&runtime.DB, user)

	if e != nil {
		glog.Errorf("failed looking up tokens: %s\n", e.Error())
		runtime.Errors = append(runtime.Errors, e)
		context.Next()
		return
	}

	for _, token := range results {
		runtime.Results = append(runtime.Results, token)
	}
	runtime.Meta.Total = total

	context.Next()
}
