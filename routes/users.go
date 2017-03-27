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

func CreateUser(runtime *net.RequestRuntime) *net.ResponseBucket {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.LogicError("bad-request")
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

		return runtime.SendErrors(errors...)
	}

	password, err := hash(body.Get("password"))

	if err != nil {
		runtime.Infof("[create user] received error hashing password: %s", err.Error())
		return runtime.FieldError("password")
	}

	email := body.Get("email")
	name := body.Get("name")

	user := models.User{Email: email, Password: password, Name: name, Username: body.Get("username")}

	usrmgr := services.UserManager{runtime.DB, runtime.Logger}

	if usrmgr.ValidPassword(body.Get("password")) != true {
		runtime.Debugf("[create user] attempt to sign up w/ invalid password: %s", body.Get("password"))
		return runtime.LogicError("invalid-password")
	}

	if ok, errors := usrmgr.ValidUser(&user); ok != true {
		runtime.Debugf("[create user] attempt to sign up w/ invalid domain: %s", email)
		return runtime.SendErrors(errors...)
	}

	if usrmgr.ValidUsername(body.Get("username")) != true {
		return runtime.LogicError("invalid-username")
	}

	if err := runtime.Create(&user).Error; err != nil {
		runtime.Errorf("[create user] unable to save: %s", err.Error())
		return runtime.ServerError()
	}

	clientmgr := services.UserClientManager{runtime.DB}
	token, err := clientmgr.Associate(&user, &runtime.Client)

	if err != nil {
		runtime.Errorf("[create user] unable to associate: %s", err.Error())
		return runtime.ServerError()
	}

	runtime.Debugf("[create user] associated user[%d] with client[%d]", user.ID, runtime.Client.ID)
	result := runtime.SendResults(1, []interface{}{user.Public()})
	result.Set("token", token.Token)
	return result
}

func UpdateUser(runtime *net.RequestRuntime) *net.ResponseBucket {
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
		return runtime.SendErrors(errors...)
	}

	return runtime.SendResults(1, []interface{}{runtime.User.Public()})
}

func FindUser(runtime *net.RequestRuntime) *net.ResponseBucket {
	blue, users := runtime.Blueprint(), []models.User{}

	count, err := blue.Apply(&users)

	if err != nil {
		runtime.Errorf("[find users] failed applying blueprint: %s", err.Error())
		return runtime.ServerError()
	}

	return runtime.SendResults(count, users)
}
