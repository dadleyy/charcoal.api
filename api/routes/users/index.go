package users

import "github.com/golang/glog"
import "github.com/kataras/iris"
import "github.com/meritoss/meritoss.api/api"
import "github.com/meritoss/meritoss.api/api/db/dal/user"

func Index(ctx *iris.Context) {
	runtime, ok := ctx.Get("runtime").(api.Runtime)

	if !ok {
		glog.Error("unable to retreive runtime from request context")
		ctx.Panic()
		return
	}

	result, err := user.Find(runtime)

	if err != nil {
		glog.Errorf("error finding users: %s", err.Error())
		ctx.Panic()
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{"results": result})
}
