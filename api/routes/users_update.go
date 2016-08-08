package routes

import "errors"
import "strconv"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/dal"
import "github.com/sizethree/meritoss.api/api/models"
import "github.com/sizethree/meritoss.api/api/middleware"


func UpdateUser(ctx *iris.Context) {
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

	// at this point we've completed our update of the user successfully. in order to 
	// give the user useful information, find and add the user to the request's jsonapi
	// bucket.
	var user models.User
	head := runtime.DB.Where("ID = ?", userid).Find(&user)

	if head.Error != nil {
		glog.Errorf("failed getting user user: %s\n", head.Error.Error())
		bucket.Errors = append(bucket.Errors, head.Error)
		ctx.Next()
		return
	}

	bucket.Results = append(bucket.Results, user)
	bucket.Meta.Total = 1

	ctx.Next()
}
