package routes

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"

func FindGameMembershipHistory(runtime *net.RequestRuntime) *net.ResponseBucket {
	results := make([]models.GameMembershipHistory, 0)
	blueprint := runtime.Blueprint(runtime.DB)

	total, err := blueprint.Apply(&results)

	if err != nil {
		runtime.Errorf("[game membership history] invalid blueprint apply: %s", err.Error())
		return runtime.ServerError()
	}

	return runtime.SendResults(total, results)
}
