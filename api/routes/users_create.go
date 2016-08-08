package routes

import "errors"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/dal"
import "github.com/sizethree/meritoss.api/api/middleware"

func CreateUser(ctx *iris.Context) {
	runtime, _ := ctx.Get("runtime").(api.Runtime)
	bucket, _ := ctx.Get("jsonapi").(*middleware.Bucket)

	var target dal.UserFacade

	if err := ctx.ReadJSON(&target); err != nil {
		bucket.Errors = append(bucket.Errors, errors.New("invalid json data for user"))
		return
	}

	target.ID = 0

	user, err := dal.CreateUser(&runtime, &target)

	if err != nil {
		bucket.Errors = append(bucket.Errors, err)
		return
	}

	bucket.Results = append(bucket.Results, user)
	bucket.Meta.Total = 1

	glog.Infof("created user %d\n", user.ID)
	ctx.Next()
}
