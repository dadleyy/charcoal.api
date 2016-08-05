package users

import "errors"
import "strconv"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/meritoss/meritoss.api/api"
import "github.com/meritoss/meritoss.api/api/dal"
import "github.com/meritoss/meritoss.api/api/middleware"


func Update(ctx *iris.Context) {
	runtime, _ := ctx.Get("runtime").(api.Runtime)
	bucket, _ := ctx.Get("jsonapi").(*middleware.Bucket)

	userid, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		bucket.Errors = append(bucket.Errors, errors.New("invalid user id"))
		ctx.Next()
		return
	}

	var updates dal.Updates

	if err := ctx.ReadJSON(&updates); err != nil {
		bucket.Errors = append(bucket.Errors, errors.New("invalid json data for user"))
		ctx.Next()
		return
	}

	if err := dal.UpdateUser(&runtime, &updates, userid); err != nil {
		glog.Errorf("failed updating user: %s\n", err.Error())
		bucket.Errors = append(bucket.Errors, err)
		ctx.Next()
		return
	}

	ctx.Next()
}
