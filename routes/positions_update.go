package routes

import "errors"
import "strconv"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/dal"
import "github.com/sizethree/meritoss.api/middleware"


func UpdatePosition(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*middleware.Runtime)

	if !ok {
		glog.Error("bad runtime found while finding clients")
		context.Panic()
		context.StopExecution()
		return
	}

	position_id, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		runtime.Error(errors.New("invalid user id"))
		context.Next()
		return
	}

	var updates struct{Location int}

	if err := context.ReadJSON(&updates); err != nil {
		runtime.Error(errors.New("invalid json data for position"))
		context.Next()
		return
	}

	facade := dal.PositionFacade{
		User: runtime.User.ID,
		Location: updates.Location,
		ID: uint(position_id),
	}

	if e := dal.UpdatePosition(&runtime.DB, &facade); e != nil {
		runtime.Error(e)
		context.Next()
		return
	}

	glog.Infof("updating position %d to %d", position_id, updates.Location)
	context.Next()
}
