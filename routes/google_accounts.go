package routes

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"

func FindGoogleAccounts(runtime *net.RequestRuntime) *net.ResponseBucket {
	var accounts []models.GoogleAccount
	blueprint := runtime.Blueprint()

	uman := services.UserManager{runtime.DB, runtime.Logger}

	// if this is not an admin user, make sure we are limiting to the current user
	if uman.IsAdmin(&runtime.User) != true {
		runtime.Debugf("user is not admin, limiting google account search to current user")
		blueprint = runtime.Blueprint(runtime.Where("user = ?", runtime.User.ID))
	}

	// limit this query to to current user only
	total, err := blueprint.Apply(&accounts)

	if err != nil {
		runtime.Errorf("[google accounts] failed blueprint lookup: %s", err.Error())
		return runtime.ServerError()
	}

	return runtime.SendResults(total, accounts)
}
