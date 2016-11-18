package routes

import "fmt"
import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"

func PrintAuth(runtime *net.RequestRuntime) error {
	runtime.AddResult(runtime.User.Public())
	return nil
}

func PrintUserRoles(runtime *net.RequestRuntime) error {
	runtime.Debugf("looking for user roles associated w/ user[%d]", runtime.User.ID)
	var maps []models.UserRoleMapping

	if err := runtime.Database().Where("user = ?", runtime.User.ID).Find(&maps).Error; err != nil {
		runtime.Debugf("failed mapping lookup: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_LOOKUP"))
	}

	if len(maps) == 0 {
		return nil
	}

	ids := make([]int64, len(maps))
	var roles []models.UserRole

	for i, mapping := range maps {
		ids[i] = int64(mapping.Role)
	}

	if err := runtime.Database().Where(ids).Find(&roles).Error; err != nil {
		runtime.Debugf("unable to associate to roles: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_ASSOCIATION"))
	}

	for _, role := range roles {
		runtime.AddResult(role.Public())
	}

	return nil
}

func PrintClientTokens(runtime *net.RequestRuntime) error {
	blueprint := runtime.Blueprint()
	var tokens []models.ClientToken
	offset := blueprint.Limit() * blueprint.Page()

	cursor := runtime.Database().Limit(blueprint.Limit()).Offset(offset).Where("client = ?", runtime.Client.ID)

	if err := cursor.Find(&tokens).Error; err != nil {
		return runtime.AddError(err)
	}

	for _, t := range tokens {
		runtime.AddResult(t)
	}

	return nil
}
