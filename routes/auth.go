package routes

import "fmt"
import "github.com/albrow/forms"
import "golang.org/x/crypto/bcrypt"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"

func PrintAuth(runtime *net.RequestRuntime) *net.ResponseBucket {
	meta := map[string]interface{}{"admin": runtime.IsAdmin()}

	return &net.ResponseBucket{
		Results: []interface{}{runtime.User.Public()},
		Meta:    meta,
	}
}

func PasswordLogin(runtime *net.RequestRuntime) *net.ResponseBucket {
	if runtime.Client.System != true {
		return runtime.LogicError("invalid-client")
	}

	body, err := forms.Parse(runtime.Request)

	if err != nil {
		runtime.Warnf("failed parsing body: %s", err.Error())
		return runtime.LogicError("invalid-body")
	}

	validator := body.Validator()

	validator.Require("email")
	validator.MatchEmail("email")

	validator.Require("password")
	validator.LengthRange("password", 6, 20)

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validator.HasErrors() == true {
		errors := make([]error, 0, len(validator.Fields()))

		for key := range validator.ErrorMap() {
			errors = append(errors, fmt.Errorf("field:%s", key))
		}

		return runtime.SendErrors(errors...)
	}

	user, token := models.User{Email: body.Get("email")}, models.ClientToken{}

	if e := runtime.Where("email = ?", user.Email).First(&user).Error; e != nil {
		runtime.Errorf("[password login] invalid login attempt: %s", e.Error())
		return runtime.LogicError("invalid-login")
	}

	if e := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Get("password"))); e != nil {
		runtime.Warnf("[password login] invalid login attempt: %s", e.Error())
		return runtime.LogicError("invalid-login")
	}

	if e := runtime.Where("user_id = ? AND client_id = ?", user.ID, runtime.Client.ID).First(&token).Error; e != nil {
		runtime.Warnf("[password login] invalid login attempt: %s", e.Error())
		return runtime.LogicError("invalid-login")
	}

	return &net.ResponseBucket{Results: []interface{}{token}}
}

func PrintUserRoles(runtime *net.RequestRuntime) *net.ResponseBucket {
	maps, roles := []models.UserRoleMapping{}, []models.UserRole{}

	runtime.Debugf("[user role lookup] looking for user roles associated w/ user[%d]", runtime.User.ID)

	if err := runtime.Where("user_id = ?", runtime.User.ID).Find(&maps).Error; err != nil {
		runtime.Warnf("failed mapping lookup: %s", err.Error())
		return runtime.ServerError()
	}

	// no user role mappings
	if len(maps) == 0 {
		return &net.ResponseBucket{}
	}

	ids := make([]int64, len(maps))

	for i, mapping := range maps {
		ids[i] = int64(mapping.RoleID)
	}

	if err := runtime.Where(ids).Find(&roles).Error; err != nil {
		runtime.Errorf("unable to associate to roles: %s", err.Error())
		return runtime.ServerError()
	}

	return &net.ResponseBucket{Results: roles}
}

func PrintClientTokens(runtime *net.RequestRuntime) *net.ResponseBucket {
	var tokens []models.ClientToken

	cursor := runtime.Where("client_id = ?", runtime.Client.ID)
	blueprint := runtime.Blueprint(cursor)

	if _, e := blueprint.Apply(&tokens); e != nil {
		runtime.Warnf("[auth client tokens] unable to lookup tokens for client %d: %s", runtime.Client.ID, e.Error())
		return runtime.ServerError()
	}

	return &net.ResponseBucket{Results: tokens}
}
