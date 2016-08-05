package users

import "errors"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/meritoss/meritoss.api/api"
import "github.com/meritoss/meritoss.api/api/dal"
import "github.com/meritoss/meritoss.api/api/models"
import "github.com/meritoss/meritoss.api/api/responses"
import "github.com/meritoss/meritoss.api/api/middleware"

func Create(ctx *iris.Context) {
	runtime, ok := ctx.Get("runtime").(api.Runtime)

	if !ok {
		responses.ServerError(ctx, "unable to lookup runtime")
		return
	}

	bucket, ok := ctx.Get("jsonapi").(*middleware.Bucket)

	if !ok {
		responses.ServerError(ctx, "unable to lookup jsonapi bucket")
		return
	}

	var target models.User

	if err := ctx.ReadJSON(&target); err != nil {
		bucket.Errors = append(bucket.Errors, errors.New("invalid json data for user"))
		return
	}

	if err := dal.CreateUser(&runtime, &target); err != nil {
		bucket.Errors = append(bucket.Errors, err)
		return
	}

	bucket.Results = append(bucket.Results, target)
	bucket.Meta.Total = 1

	glog.Infof("created user %d\n", target.ID)
	ctx.Next()
}
