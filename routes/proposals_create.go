package routes

import "errors"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/dal"
import "github.com/sizethree/meritoss.api/middleware"

func CreateProposal(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*middleware.Runtime)

	if !ok {
		glog.Error("bad runtime found while finding clients")
		context.Panic()
		context.StopExecution()
		return
	}

	var target dal.ProposalFacade

	if err := context.ReadJSON(&target); err != nil {
		runtime.Error(errors.New("invalid json data for user"))
		return
	}

	proposal, err := dal.CreateProposal(&runtime.DB, &target)

	if err != nil {
		runtime.Error(err)
		return
	}

	runtime.Result(proposal)
	runtime.Meta("total", 1)

	glog.Infof("created proposal %d\n", proposal.ID)
	context.Next()
}
