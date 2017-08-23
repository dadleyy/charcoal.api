package routes

import "fmt"
import "strconv"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/charcoal/net"
import "github.com/dadleyy/charcoal.api/charcoal/models"
import "github.com/dadleyy/charcoal.api/charcoal/services"

func DestroyUserRoleMapping(runtime *net.RequestRuntime) *net.ResponseBucket {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.LogicError("invalid-id")
	}

	var mapping models.UserRoleMapping

	if err := runtime.Where("id = ?", id).First(&mapping).Error; err != nil {
		runtime.Warnf("[destroy mapping] unable to find mapping: %s", err.Error())
		return runtime.LogicError("not-found")
	}

	if err := runtime.Unscoped().Delete(&mapping).Error; err != nil {
		runtime.Debugf("[del mapping] unable to delete role mapping: %s", err.Error())
		return runtime.ServerError()
	}

	return nil
}

func CreateUserRoleMapping(runtime *net.RequestRuntime) *net.ResponseBucket {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.LogicError("invalid-body")
	}

	validator := body.Validator()
	validator.Require("user")
	validator.Require("role")

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validator.HasErrors() == true {
		for _, m := range validator.Messages() {
			runtime.AddError(fmt.Errorf(m))
		}

		return nil
	}

	user := body.Get("user_id")
	role := body.Get("role_id")

	uid, err := strconv.Atoi(user)

	if err != nil {
		return runtime.LogicError("missing-user")
	}

	rid, err := strconv.Atoi(role)

	if err != nil {
		return runtime.LogicError("missing-role")
	}

	mapping, duplicate := models.UserRoleMapping{RoleID: uint(rid), UserID: uint(uid)}, 0

	cursor := runtime.Model(&mapping).Where("user_id = ? AND role_id = ?", uid, rid)

	if _ = cursor.Count(&duplicate); duplicate >= 1 {
		return runtime.LogicError("duplicate")
	}

	if err := runtime.Create(&mapping).Error; err != nil {
		runtime.Errorf("[create mapping] failed created: %s", err.Error())
		return runtime.ServerError()
	}

	return runtime.SendResults(1, []models.UserRoleMapping{mapping})
}

func FindUserRoleMappings(runtime *net.RequestRuntime) *net.ResponseBucket {
	var maps []models.UserRoleMapping
	blueprint := runtime.Blueprint()

	uman := services.UserManager{runtime.DB, runtime.Logger}

	// if this is not an admin user, make sure we are limiting to the current user
	if uman.IsAdmin(&runtime.User) != true {
		runtime.Debugf("user is not admin, limiting role maps search to current user")
		blueprint = runtime.Blueprint(runtime.Where("user_id = ?", runtime.User.ID))
	}

	// limit this query to to current user only
	total, err := blueprint.Apply(&maps)

	if err != nil {
		runtime.Errorf("[find mappings] unable to apply blueprint: %s", err.Error())
		return runtime.ServerError()
	}

	return runtime.SendResults(total, maps)
}