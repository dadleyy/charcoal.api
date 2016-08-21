package middleware

import "time"
import "github.com/golang/glog"
import "github.com/dadleyy/iris"

func Logger(ctx *iris.Context) {
	start := time.Now()

	defer func() {
		since := time.Since(start).Seconds()
		glog.Infof("\"%s %s\" took %f seconds\n", ctx.MethodString(), ctx.PathString(), since)
	}()

	ctx.Next()
}

