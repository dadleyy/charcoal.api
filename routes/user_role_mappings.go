package routes

import "fmt"
import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/services"

func FindUserRoleMappings(runtime *net.RequestRuntime) error {
	var maps []models.UserRoleMapping
	blueprint := runtime.Blueprint()

	uman := services.UserManager{runtime.Database()}

	// if this is not an admin user, make sure we are limiting to the current user
	if uman.IsAdmin(&runtime.User) != true {
		runtime.Debugf("user is not admin, limiting role maps search to current user")
		blueprint.Filter("filter[user]", fmt.Sprintf("eq(%d)", runtime.User.ID))
	}

	// limit this query to to current user only
	total, err := blueprint.Apply(&maps, runtime.Database())

	if err != nil {
		return err
	}

	for _, item := range maps {
		runtime.AddResult(item)
	}

	runtime.SetMeta("toal", total)

	return nil
}
