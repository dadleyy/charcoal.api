package routes

import "errors"
import "strconv"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/dal"
import "github.com/sizethree/meritoss.api/models"
import "github.com/sizethree/meritoss.api/middleware"

// UpdateUser
//
// request callback for PATCH /users/:id
//
// attempts to load in a set of updates per data access layer definition and use the
// dal.UpdateUser function to apply those updates to the persistance
func UpdateUser(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*middleware.Runtime)

	if !ok {
		glog.Error("bad runtime")
		context.Panic()
		context.StopExecution()
		return
	}

	userid, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		runtime.Error(errors.New("invalid user id"))
		context.Next()
		return
	}

	if userid != int(runtime.User.ID) {
		runtime.Error(errors.New("attempt to update different user"))
		context.Next()
		return
	}

	var updates dal.Updates

	if err := context.ReadJSON(&updates); err != nil {
		runtime.Error(errors.New("invalid json data for user"))
		context.Next()
		return
	}

	if e := dal.UpdateUser(&runtime.DB, &updates, userid); e != nil {
		glog.Errorf("unable to update user: %s\n", e.Error())
		runtime.Error(e)
		context.Next()
		return
	}

	var user models.User

	if e := runtime.DB.Where("id = ?", userid).First(&user).Error; e != nil {
		runtime.Error(e)
		context.Next()
		return
	}

	runtime.Result(user)
	runtime.Meta("total", 1)

	context.Next()
}
