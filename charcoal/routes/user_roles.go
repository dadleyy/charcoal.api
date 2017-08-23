package routes

import "github.com/dadleyy/charcoal.api/charcoal/net"
import "github.com/dadleyy/charcoal.api/charcoal/models"

func FindRoles(runtime *net.RequestRuntime) *net.ResponseBucket {
	blueprint, roles := runtime.Blueprint(), []models.UserRole{}
	total, err := blueprint.Apply(&roles)

	if err != nil {
		runtime.Warnf("[find roles] badd lookup: %s", err.Error())
		return runtime.LogicError("bad-request")
	}

	return runtime.SendResults(total, roles)
}
