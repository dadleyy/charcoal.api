package routes

import "fmt"

import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/services"

func FindClientAdmins(runtime *net.RequestRuntime) error {
	var results []models.ClientAdmin
	blueprint := runtime.Blueprint()

	uman := services.UserManager{runtime.Database()}

	if uman.IsAdmin(&runtime.User) != true {
		runtime.Debugf("user is not admin, limiting query to client[%d]", runtime.Client.ID)
		err := blueprint.Filter("filter[client]", fmt.Sprintf("eq(%d)", runtime.Client.ID))

		if err != nil {
			runtime.Debugf("filter problem: %s", err.Error())
			return runtime.AddError(fmt.Errorf("PROBLEM"))
		}

		// make sure user is even able to see this client's admins
	}

	total, err := blueprint.Apply(&results, runtime.Database())

	if err != nil {
		runtime.Debugf("BAD_LOOKUP: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_QUERY"))
	}

	for _, item := range results {
		runtime.AddResult(item)
	}

	runtime.SetMeta("total", total)

	return nil
}
