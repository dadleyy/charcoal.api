package middleware

import "errors"
import "github.com/golang/glog"
import "github.com/kataras/iris"
import "github.com/sizethree/meritoss.api/api"

func RequireAuth(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*api.Runtime)

	if !ok || runtime.User.ID < 1 {
		glog.Errorf("[authentication error] user authentication required for %s", context.Path())
		runtime.Error(errors.New("not found"))
		runtime.Render(context)
		return
	}

	context.Next()
}

