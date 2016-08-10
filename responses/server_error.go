package responses

import "github.com/golang/glog"
import "github.com/kataras/iris"

func ServerError(ctx *iris.Context, message string) {
	glog.Error(message)
	ctx.Panic()
	ctx.StopExecution()
}
