package api

import "github.com/kataras/iris"
import "github.com/sizethree/meritoss.api/api/db"
import "github.com/sizethree/meritoss.api/api/models"

type Runtime struct {
	*Bucket
	DB db.Client
	User models.User
}

func (runtime *Runtime) Finish(context *iris.Context) {
	runtime.DB.Close()
	runtime.Render(context)
	context.StopExecution()
}
