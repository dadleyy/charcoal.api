package routes
/*

import "fmt"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/services"

const ERR_BAD_SESSION = "BAD_SESSION"

func PrintAuth(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Runtime)

	if ok != true {
		return fmt.Errorf(ERR_BAD_SESSION)
	}

	runtime.AddResult(&runtime.User)

	uman := services.UserManager{runtime.DB}
	runtime.AddMeta("admin", uman.IsAdmin(&runtime.User))

	return nil
}

func PrintUserRoles(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Runtime)

	if ok != true {
		return fmt.Errorf(ERR_BAD_SESSION)
	}

	runtime.Logger().Debugf("looking for user roles associated w/ user[%d]", runtime.User.ID)

	var maps []models.UserRoleMapping

	if err := runtime.DB.Where("user = ?", runtime.User.ID).Find(&maps).Error; err != nil{
		runtime.Logger().Debugf("failed mapping lookup: %s", err.Error())
		return fmt.Errorf("BAD_LOOKUP")
	}

	if len(maps) == 0 {
		return nil
	}

	ids := make([]int64, len(maps))
	var roles []models.UserRole

	for i, mapping := range maps {
		ids[i] = int64(mapping.Role)
	}

	if err := runtime.DB.Where(ids).Find(&roles).Error; err != nil {
		runtime.Logger().Debugf("unable to associate to roles: %s", err.Error())
		return fmt.Errorf("BAD_ASSOCIATION")
	}

	for _, role := range roles {
		runtime.AddResult(role)
	}

	return nil
}

func PrintClientTokens(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Runtime)

	if ok != true {
		return fmt.Errorf(ERR_BAD_SESSION)
	}

	if runtime.Client.ID == 0 {
		return fmt.Errorf("BAD_CLIENT")
	}

	blueprint := runtime.Blueprint()
	var tokens []models.ClientToken

	blueprint.Filter("filter[client]", fmt.Sprint("eq(%d)", runtime.Client.ID))

	total, err := blueprint.Apply(&tokens, runtime.DB)

	if err != nil {
		runtime.Logger().Debugf("unable to apply client tokens: %s", err.Error())
		return fmt.Errorf("NO_TOKENS")
	}

	for _, token := range tokens {
		runtime.AddResult(token)
	}

	runtime.AddMeta("total", total)

	return nil
}
*/
