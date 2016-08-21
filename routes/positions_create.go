package routes

import "errors"
import "github.com/golang/glog"
import "github.com/dadleyy/iris"

import "github.com/sizethree/miritos.api/dal"
import "github.com/sizethree/miritos.api/middleware"

func CreatePosition(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*middleware.Runtime)

	if !ok {
		glog.Error("bad runtime found while finding clients")
		context.Panic()
		context.StopExecution()
		return
	}

	var target dal.PositionFacade

	if err := context.ReadJSON(&target); err != nil {
		runtime.Error(errors.New("invalid json data for user"))
		context.Next()
		return
	}

	target.User = runtime.User.ID

	position, err := dal.CreatePosition(&runtime.DB, &target)

	if err != nil {
		glog.Errorf("unable to create position: %s", err.Error())
		runtime.Error(err)
		context.Next()
		return
	}

	runtime.Result(position)
	glog.Infof("creating position!")
	context.Next()
}
