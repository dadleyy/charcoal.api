package routes

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"

func FindRoles(runtime *net.RequestRuntime) *net.ResponseBucket {
	blueprint, roles := runtime.Blueprint(), []models.UserRole{}
	total, err := blueprint.Apply(&roles)

	if err != nil {
		runtime.Warnf("[find roles] badd lookup: %s", err.Error())
		return runtime.LogicError("bad-request")
	}

	return runtime.SendResults(total, roles)
}
