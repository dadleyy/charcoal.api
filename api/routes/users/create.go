package users

import "errors"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/meritoss/meritoss.api/api"
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


	var user models.User

	if err := ctx.ReadJSON(&user); err != nil {
		bucket.Errors = append(bucket.Errors, errors.New("invalid json data for user"))
		return
	}

	if len(user.Name) < 2 {
		bucket.Errors = append(bucket.Errors, errors.New("user name must be at least 2 characters long"))
		return
	}

	if err := runtime.DB.Save(&user).Error; err != nil {
		glog.Errorf("error saving user: %s\n", err.Error())
		bucket.Errors = append(bucket.Errors, errors.New("unable to save user"))
		return
	}

	bucket.Results = append(bucket.Results, user)
	bucket.Meta.Total = 1

	glog.Infof("adding user %d to bucket. count: %d\n", user.ID, len(bucket.Results))


	ctx.Next()
}
