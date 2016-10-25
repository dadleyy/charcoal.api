package routes

import "fmt"
import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/services"

func FindGoogleAccounts(runtime *net.RequestRuntime) error {
	var accounts []models.GoogleAccount
	blueprint := runtime.Blueprint()

	uman := services.UserManager{runtime.Database()}

	// if this is not an admin user, make sure we are limiting to the current user
	if uman.IsAdmin(&runtime.User) != true {
		runtime.Debugf("user is not admin, limiting google account search to current user")
		blueprint.Filter("filter[user]", fmt.Sprintf("eq(%d)", runtime.User.ID))
	}

	// limit this query to to current user only
	total, err := blueprint.Apply(&accounts, runtime.Database())

	if err != nil {
		return err
	}

	for _, account := range accounts {
		runtime.AddResult(account)
	}

	runtime.SetMeta("toal", total)

	return nil
}
