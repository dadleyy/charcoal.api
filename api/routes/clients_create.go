package routes

import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/dal"
import "github.com/sizethree/meritoss.api/api/middleware"

func CreateClient(context *iris.Context) {
	runtime, _ := context.Get("runtime").(api.Runtime)
	bucket, _ := context.Get("jsonapi").(*middleware.Bucket)

	var target dal.ClientFacade

	if e := context.ReadJSON(&target); e != nil {
		bucket.Errors = append(bucket.Errors, e)
		return
	}

	client, e := dal.CreateClient(&runtime, &target);

	if e != nil {
		bucket.Errors = append(bucket.Errors, e)
		return
	}

	bucket.Results = append(bucket.Results, client)
	bucket.Meta.Total = 1

	glog.Infof("created client %d\n", client.ID)
	context.Next()
}
