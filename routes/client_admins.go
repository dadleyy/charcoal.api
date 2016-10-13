package routes
/*

import "fmt"
import "github.com/labstack/echo"

import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/services"

func FindClientAdmins(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Runtime)

	var results []models.ClientAdmin
	blueprint := runtime.Blueprint()

	uman := services.UserManager{runtime.DB}

	if uman.IsAdmin(&runtime.User) != true {
		runtime.Logger().Debugf("user is not admin, limiting query to client[%d]", runtime.Client.ID)
		err := blueprint.Filter("filter[client]", fmt.Sprintf("eq(%d)", runtime.Client.ID))

		if err != nil {
			runtime.Logger().Debugf("filter problem: %s", err.Error())
			return fmt.Errorf("PROBLEM")
		}

		// make sure user is even able to see this client's admins
	}

	total, err := blueprint.Apply(&results, runtime.DB)

	if err != nil {
		runtime.Logger().Debugf("BAD_LOOKUP: %s", err.Error())
		return fmt.Errorf("BAD_QUERY")
	}

	for _, item := range results {
		runtime.AddResult(item)
	}

	runtime.AddMeta("total", total)

	return nil
}
*/
