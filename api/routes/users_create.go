package routes

import "errors"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/dal"

func CreateUser(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*api.Runtime)

	if !ok {
		glog.Error("bad runtime!")
		context.Panic()
		context.StopExecution()
		return
	}

	var target dal.UserFacade

	if err := context.ReadJSON(&target); err != nil {
		runtime.Errors = append(runtime.Errors, errors.New("invalid json data for user"))
		return
	}

	user, err := dal.CreateUser(&runtime.DB, &target)

	if err != nil {
		runtime.Errors = append(runtime.Errors, err)
		return
	}

	runtime.Results = append(runtime.Results, user)
	runtime.Meta.Total = 1

	glog.Infof("created user %d\n", user.ID)
	context.Next()
}
