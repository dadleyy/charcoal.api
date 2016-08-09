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
		runtime.Error(errors.New("invalid json data for user"))
		context.Next()
		return
	}

	user, err := dal.CreateUser(&runtime.DB, &target)

	if err != nil {
		runtime.Error(err)
		return
	}

	runtime.Result(user)
	runtime.Meta("total", 1)

	glog.Infof("created user %d\n", user.ID)
	context.Next()
}
