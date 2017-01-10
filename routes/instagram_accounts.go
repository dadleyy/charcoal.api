package routes

/*
import "github.com/sizethree/miritos.api/services"
*/

import "fmt"
import "github.com/albrow/forms"
import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"

func CreateInstagramAccount(runtime *net.RequestRuntime) error {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.AddError(err)
	}

	validator := body.Validator()
	validator.Require("instagram_id")
	validator.Require("username")

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validator.HasErrors() == true {
		for _, m := range validator.Messages() {
			runtime.AddError(fmt.Errorf(m))
		}

		return nil
	}

	username := body.Get("username")
	gramid := body.Get("instagram_id")

	runtime.Debugf("checking insta-user[%s]: %s", gramid, username)

	var count uint = 0
	if runtime.Cursor(&models.InstagramAccount{}).Where("instagram_id = ?", gramid).Count(&count); count >= 1 {
		return runtime.AddError(fmt.Errorf("DUPLICATE"))
	}

	account := models.InstagramAccount{InstagramID: gramid, Username: username}

	if err := runtime.Database().Create(&account).Error; err != nil {
		return runtime.AddError(err)
	}

	runtime.AddResult(account.Public())

	return nil
}

func FindInstagramAccounts(runtime *net.RequestRuntime) error {
	blueprint := runtime.Blueprint()
	var accounts []models.InstagramAccount

	total, err := blueprint.Apply(&accounts, runtime.Database())

	if err != nil {
		runtime.Debugf("unable to query accounts: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_QUERY"))
	}

	for _, account := range accounts {
		runtime.AddResult(account.Public())
	}

	runtime.SetMeta("total", total)

	return nil
}
