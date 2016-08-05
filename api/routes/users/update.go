package users

import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/meritoss/meritoss.api/api"

func Update(ctx *iris.Context) {
	_, ok := ctx.Get("runtime").(api.Runtime)
	glog.Infof("created user %t\n", ok)
	ctx.Next()
}
