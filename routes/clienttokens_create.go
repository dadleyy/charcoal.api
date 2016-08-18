package routes

import "errors"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/miritos.api/dal"
import "github.com/sizethree/miritos.api/middleware"

func CreateClientToken(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*middleware.Runtime)

	if !ok {
		glog.Error("bad runtime")
		context.Panic()
		context.StopExecution()
		return
	}

	var target dal.ClientTokenFacade

	if err := context.ReadJSON(&target); err != nil {
		runtime.Error(errors.New("invalid json data for user"))
		context.Next()
		return
	}

	if target.User != runtime.User.ID {
		runtime.Error(errors.New("unauthorized user"))
		context.Next()
		return
	}

	target.Referrer = runtime.Client

	glog.Infof("making token for user %d and client %d from client %d\n", target.User, target.Client, target.Referrer.ID)

	result, e := dal.CreateClientToken(&runtime.DB, &target)

	if e != nil {
		runtime.Error(e)
		context.Next()
		return
	}

	runtime.Result(result)
	context.Next()
}
