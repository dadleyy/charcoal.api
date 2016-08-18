package routes

import "errors"
import "strconv"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/miritos.api/dal"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/middleware"

func UpdateProposal(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*middleware.Runtime)

	if !ok {
		glog.Error("bad runtime found while finding clients")
		context.Panic()
		context.StopExecution()
		return
	}

	propid, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		runtime.Error(errors.New("invalid proposal id"))
		context.Next()
		return
	}

	var updates dal.Updates

	if err := context.ReadJSON(&updates); err != nil {
		runtime.Error(errors.New("invalid json data for proposal"))
		context.Next()
		return
	}

	if e := dal.UpdateProposal(&runtime.DB, &updates, propid, runtime.User.ID); e != nil {
		glog.Errorf("unable to update proposal: %s\n", e.Error())
		runtime.Error(e)
		context.Next()
		return
	}

	var prop models.Proposal

	if e := runtime.DB.Where("id = ?", propid).First(&prop).Error; e != nil {
		runtime.Error(e)
		context.Next()
		return
	}

	runtime.Result(prop)
	runtime.Meta("total", 1)

	glog.Infof("updated proposal %d", propid)
	context.Next()
}
