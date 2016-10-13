package routes
/*

import "fmt"
import "github.com/labstack/echo"

import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

func FindRoles(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Runtime)

	blueprint := runtime.Blueprint()

	var roles []models.UserRole

	total, err := blueprint.Apply(&roles, runtime.DB)

	if err != nil {
		runtime.Logger().Debugf("ERR_BAD_ROLE_LOOKUP: %s", err.Error())
		return fmt.Errorf("BAD_QUERY")
	}

	for _, role := range roles {
		runtime.AddResult(role)
	}

	runtime.AddMeta("total", total)

	return nil
}
*/
