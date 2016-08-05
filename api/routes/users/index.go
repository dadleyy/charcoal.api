package users

import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/dal"
import "github.com/sizethree/meritoss.api/api/middleware"

func Index(ctx *iris.Context) {
	runtime, _ := ctx.Get("runtime").(api.Runtime)
	bucket, _ := ctx.Get("jsonapi").(*middleware.Bucket)
	blueprint, _ := ctx.Get("blueprint").(middleware.Blueprint)

	result, total, err := dal.FindUser(runtime, blueprint)

	if err != nil {
		bucket.Errors = append(bucket.Errors, err)
		return
	}

	for _, user := range result {
		bucket.Results = append(bucket.Results, user)
	}

	bucket.Meta.Total = total

	ctx.Next()
}
