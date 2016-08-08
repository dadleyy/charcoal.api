package routes

import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/dal"
import "github.com/sizethree/meritoss.api/api/middleware"

func FindProposals(context *iris.Context) {
	runtime, _ := context.Get("runtime").(api.Runtime)
	bucket, _ := context.Get("jsonapi").(*middleware.Bucket)
	blueprint, _ := context.Get("blueprint").(*api.Blueprint)

	result, total, err := dal.FindProposals(&runtime, blueprint)

	if err != nil {
		bucket.Errors = append(bucket.Errors, err)
		return
	}

	for _, proposal := range result {
		bucket.Results = append(bucket.Results, proposal)
	}

	bucket.Meta.Total = total

	glog.Infof("searching for proposals\n")
	context.Next()
}
