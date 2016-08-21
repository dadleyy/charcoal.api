package routes

import "github.com/golang/glog"
import "github.com/dadleyy/iris"

func FindPositions(context *iris.Context) {
	glog.Infof("finding positions!")
	context.Next()
}
