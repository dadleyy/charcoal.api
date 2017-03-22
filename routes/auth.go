package routes

import "fmt"
import "github.com/albrow/forms"
import "golang.org/x/crypto/bcrypt"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"

func PrintAuth(runtime *net.RequestRuntime) error {
	runtime.AddResult(runtime.User.Public())
	runtime.SetMeta("admin", runtime.IsAdmin())
	return nil
}

func PasswordLogin(runtime *net.RequestRuntime) error {
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

		return runtime.AddError(errors...)
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

	runtime.AddResult(token)
	return nil
}

func PrintUserRoles(runtime *net.RequestRuntime) error {
	runtime.Debugf("looking for user roles associated w/ user[%d]", runtime.User.ID)
	var maps []models.UserRoleMapping

	if err := runtime.Where("user = ?", runtime.User.ID).Find(&maps).Error; err != nil {
		runtime.Warnf("failed mapping lookup: %s", err.Error())
		return runtime.ServerError()
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
		runtime.Warnf("unable to associate to roles: %s", err.Error())
		return runtime.ServerError()
	}

	for _, role := range roles {
		runtime.AddResult(role.Public())
	}

	return nil
}

func PrintClientTokens(runtime *net.RequestRuntime) error {
	var tokens []models.ClientToken

	cursor := runtime.Where("client = ?", runtime.Client.ID)
	blueprint := runtime.Blueprint(cursor)

	if _, err := blueprint.Apply(&tokens); err != nil {
		runtime.Warnf("unable to lookup tokens for client %d: %s", runtime.Client.ID, "")
		return runtime.ServerError()
	}

	for _, t := range tokens {
		runtime.AddResult(t)
	}

	return nil
}
