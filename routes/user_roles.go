package routes

import "fmt"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"

func FindRoles(runtime *net.RequestRuntime) error {
	blueprint := runtime.Blueprint()
	var roles []models.UserRole

	total, err := blueprint.Apply(&roles, runtime.Database())

	if err != nil {
		runtime.Debugf("ERR_BAD_ROLE_LOOKUP: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_QUERY"))
	}

	for _, role := range roles {
		runtime.AddResult(role)
	}

	runtime.SetMeta("total", total)

	return nil
}
