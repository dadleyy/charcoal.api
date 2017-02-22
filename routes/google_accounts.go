package routes

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"

func FindGoogleAccounts(runtime *net.RequestRuntime) error {
	var accounts []models.GoogleAccount
	blueprint := runtime.Blueprint()

	uman := services.UserManager{runtime.DB}

	// if this is not an admin user, make sure we are limiting to the current user
	if uman.IsAdmin(&runtime.User) != true {
		runtime.Debugf("user is not admin, limiting google account search to current user")
		blueprint = runtime.Blueprint(runtime.Where("user = ?", runtime.User.ID))
	}

	// limit this query to to current user only
	total, err := blueprint.Apply(&accounts)

	if err != nil {
		return err
	}

	for _, account := range accounts {
		runtime.AddResult(account)
	}

	runtime.SetMeta("toal", total)

	return nil
}
