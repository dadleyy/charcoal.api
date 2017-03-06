package routes

import "fmt"
import "strconv"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"

func DestroyUserRoleMapping(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.AddError(fmt.Errorf("BAD_PHOTO_ID"))
	}

	var mapping models.UserRoleMapping

	if err := runtime.Where("id = ?", id).First(&mapping).Error; err != nil {
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	if err := runtime.Unscoped().Delete(&mapping).Error; err != nil {
		runtime.Debugf("unable to delete role mapping: %s", err.Error())
		return runtime.AddError(fmt.Errorf("CANT_DELETE"))
	}

	return nil
}

func CreateUserRoleMapping(runtime *net.RequestRuntime) error {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return err
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

	user := body.Get("user")
	role := body.Get("role")

	uid, err := strconv.Atoi(user)

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_USER"))
	}

	rid, err := strconv.Atoi(role)

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_ROLE"))
	}

	mapping := models.UserRoleMapping{Role: uint(rid), User: uint(uid)}
	duplicate := 0

	cursor := runtime.Model(&mapping).Where("user = ? AND role = ?", uid, rid)

	if _ = cursor.Count(&duplicate); duplicate >= 1 {
		return runtime.AddError(fmt.Errorf("MAPPING_EXISTS"))
	}

	if err := runtime.Create(&mapping).Error; err != nil {
		return runtime.AddError(fmt.Errorf("BAD_CREATE"))
	}

	runtime.AddResult(mapping)

	return nil
}

func FindUserRoleMappings(runtime *net.RequestRuntime) error {
	var maps []models.UserRoleMapping
	blueprint := runtime.Blueprint()

	uman := services.UserManager{runtime.DB, runtime.Logger}

	// if this is not an admin user, make sure we are limiting to the current user
	if uman.IsAdmin(&runtime.User) != true {
		runtime.Debugf("user is not admin, limiting role maps search to current user")
		blueprint = runtime.Blueprint(runtime.Where("user = ?", runtime.User.ID))
	}

	// limit this query to to current user only
	total, err := blueprint.Apply(&maps)

	if err != nil {
		return err
	}

	for _, item := range maps {
		runtime.AddResult(item)
	}

	runtime.SetMeta("toal", total)

	return nil
}
