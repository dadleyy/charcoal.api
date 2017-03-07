package routes

import "fmt"
import "github.com/albrow/forms"
import "golang.org/x/crypto/bcrypt"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"

func hash(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(result), nil
}

func CreateUser(runtime *net.RequestRuntime) error {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.AddError(err)
	}

	validator := body.Validator()

	validator.Require("email")
	validator.MatchEmail("email")

	validator.Require("password")
	validator.LengthRange("password", 6, 20)

	validator.Require("name")
	validator.LengthRange("password", 2, 100)

	validator.Require("username")
	validator.LengthRange("username", 2, 100)

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validator.HasErrors() == true {
		errors := make([]error, 0, len(validator.Fields()))

		for key := range validator.ErrorMap() {
			errors = append(errors, fmt.Errorf("field:%s", key))
		}

		return runtime.AddError(errors...)
	}

	password, err := hash(body.Get("password"))

	if err != nil {
		runtime.Infof("received error hashing password: %s", err.Error())
		return runtime.FieldError("password")
	}

	email := body.Get("email")
	name := body.Get("name")

	user := models.User{Email: email, Password: password, Name: name, Username: body.Get("username")}

	usrmgr := services.UserManager{runtime.DB, runtime.Logger}

	if usrmgr.ValidPassword(body.Get("password")) != true {
		runtime.Debugf("attempt to sign up w/ invalid password: %s", body.Get("password"))
		return runtime.LogicError("invalid-password")
	}

	if ok, errors := usrmgr.ValidUser(&user); ok != true {
		runtime.Debugf("attempt to sign up w/ invalid domain: %s", email)
		return runtime.AddError(errors...)
	}

	if usrmgr.ValidUsername(body.Get("username")) != true {
		return runtime.LogicError("invalid-username")
	}

	if err := runtime.Create(&user).Error; err != nil {
		runtime.Debugf("unable to save: %s", err.Error())
		return runtime.ServerError()
	}

	clientmgr := services.UserClientManager{runtime.DB}
	token, err := clientmgr.Associate(&user, &runtime.Client)

	if err != nil {
		runtime.Debugf("unable to associate: %s", err.Error())
		return runtime.ServerError()
	}

	runtime.Debugf("associated user[%d] with client[%d]", user.ID, runtime.Client.ID)
	runtime.AddResult(user.Public())
	runtime.SetMeta("token", token.Token)

	return nil
}

func UpdateUser(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.LogicError("bad-user-id")
	}

	if runtime.User.ID != uint(id) {
		return runtime.LogicError("not-authorized")
	}

	body, err := forms.Parse(runtime.Request)

	if err != nil {
		runtime.Warnf("error parsing body: %s", err.Error())
		return runtime.LogicError("bad-body")
	}

	usrmgr := services.UserManager{runtime.DB, runtime.Logger}

	if errors := usrmgr.ApplyUpdates(&runtime.User, body.Values); len(errors) >= 1 {
		runtime.Debugf("update to user[%d] failed - %v", id, errors)
		return runtime.AddError(errors...)
	}

	runtime.AddResult(runtime.User.Public())
	return nil
}

func FindUser(runtime *net.RequestRuntime) error {
	var users []models.User
	blue := runtime.Blueprint()

	count, err := blue.Apply(&users)

	if err != nil {
		runtime.Debugf("failed applying blueprint: %s", err.Error())
		return runtime.ServerError()
	}

	runtime.SetTotal(count)

	for _, u := range users {
		runtime.AddResult(u.Public())
	}

	return nil
}
