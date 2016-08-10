package routes

import "errors"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/dal"
import "github.com/sizethree/meritoss.api/models"
import "github.com/sizethree/meritoss.api/middleware"

func CreateUser(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*middleware.Runtime)

	if !ok {
		glog.Error("bad runtime!")
		context.Panic()
		context.StopExecution()
		return
	}

	target := dal.UserFacade{ReferrerClient: runtime.Client}

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

	// now that the user is created, try finding the client token that is
	// associated with the client that created the user and the user that
	// was created so that we can add that to the meta data.
	var token models.ClientToken
	runtime.DB.Where("client = ? AND user = ?", runtime.Client.ID, user.ID).First(&token)

	runtime.Result(user)
	runtime.Meta("total", 1)
	runtime.Meta("token", token.Token)

	glog.Infof("created user %d\n", user.ID)
	context.Next()
}
