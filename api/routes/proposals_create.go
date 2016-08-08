package routes

import "errors"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/dal"
import "github.com/sizethree/meritoss.api/api/middleware"

func CreateProposal(ctx *iris.Context) {
	runtime, _ := ctx.Get("runtime").(api.Runtime)
	bucket, _ := ctx.Get("jsonapi").(*middleware.Bucket)

	var target dal.ProposalFacade

	if err := ctx.ReadJSON(&target); err != nil {
		bucket.Errors = append(bucket.Errors, errors.New("invalid json data for user"))
		return
	}

	target.ID = 0

	proposal, err := dal.CreateProposal(&runtime, &target);

	if err != nil {
		bucket.Errors = append(bucket.Errors, err)
		return
	}

	bucket.Results = append(bucket.Results, proposal)
	bucket.Meta.Total = 1

	glog.Infof("created proposal %d\n", proposal.ID)
	ctx.Next()
}
