package middleware

import "errors"
import "github.com/golang/glog"
import "github.com/kataras/iris"

func RequireAuth(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*Runtime)

	if !ok || runtime.User.ID < 1 {
		glog.Errorf("[authentication error] user authentication required for %s", context.Path())
		runtime.Error(errors.New("not found"))
		runtime.Render(context)
		return
	}

	context.Next()
}

