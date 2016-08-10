package routes

import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/dal"
import "github.com/sizethree/meritoss.api/middleware"

func FindClientTokens(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*middleware.Runtime)

	if !ok {
		glog.Error("bad runtime")
		context.Panic()
		context.StopExecution()
		return
	}

	user := runtime.User.ID
	glog.Infof("looking up client tokens for user %d\n", user)

	results, total, e := dal.ClientTokensForUser(&runtime.DB, user)

	if e != nil {
		glog.Errorf("failed looking up tokens: %s\n", e.Error())
		runtime.Error(e)
		context.Next()
		return
	}

	// add each token to the result
	for _, token := range results {
		runtime.Result(token)
	}

	runtime.Meta("total", total)

	context.Next()
}
