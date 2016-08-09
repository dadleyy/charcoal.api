package routes

import "errors"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/dal"

func CreateClientToken(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*api.Runtime)

	if !ok {
		glog.Error("bad runtime")
		context.Panic()
		context.StopExecution()
		return
	}

	var target dal.ClientTokenFacade

	if err := context.ReadJSON(&target); err != nil {
		runtime.Errors = append(runtime.Errors, errors.New("invalid json data for user"))
		context.Next()
		return
	}

	if target.User != runtime.User.ID {
		runtime.Errors = append(runtime.Errors, errors.New("unauthorized user"))
		context.Next()
		return
	}

	target.Referrer = runtime.Client

	glog.Infof("making token for user %d and client %d from client %d\n", target.User, target.Client, target.Referrer.ID)

	result, e := dal.CreateClientToken(&runtime.DB, &target)

	if e != nil {
		runtime.Errors = append(runtime.Errors, e)
		context.Next()
		return
	}

	runtime.Results = append(runtime.Results, result)
	context.Next()
}
