package routes

import "fmt"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"

func PrintAuth(runtime *net.RequestRuntime) error {
	runtime.AddResult(runtime.User.Public())
	runtime.SetMeta("admin", runtime.IsAdmin())
	return nil
}

func PrintUserRoles(runtime *net.RequestRuntime) error {
	runtime.Debugf("looking for user roles associated w/ user[%d]", runtime.User.ID)
	var maps []models.UserRoleMapping

	if err := runtime.Where("user = ?", runtime.User.ID).Find(&maps).Error; err != nil {
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

	if err := runtime.Where(ids).Find(&roles).Error; err != nil {
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

	cursor := runtime.Limit(blueprint.Limit()).Offset(offset).Where("client = ?", runtime.Client.ID)

	if err := cursor.Find(&tokens).Error; err != nil {
		return runtime.AddError(err)
	}

	for _, t := range tokens {
		runtime.AddResult(t)
	}

	return nil
}
