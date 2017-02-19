package routes

import "fmt"

import "strconv"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"

func DeleteClientAdmin(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.AddError(fmt.Errorf("BAD_ID"))
	}

	var record models.ClientAdmin

	if err := runtime.Database().First(&record, id).Error; err != nil {
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	authorized := runtime.IsAdmin()

	// if the user is not a system admin, check to see if they are an admin of the client
	if authorized != true {
		count := 0
		cursor := runtime.Cursor(&models.ClientAdmin{}).Where("client = ? AND user = ?", record.Client, runtime.User.ID)

		if err := cursor.Count(&count).Error; err != nil || count == 0 {
			message := "unauthorized attempt to remove client admin user[%d] record[%d]: %v"
			runtime.Debugf(message, runtime.User.ID, id, err)
			return runtime.AddError(fmt.Errorf("UNAUTHORIZED"))
		}
	}

	if record.User == runtime.User.ID {
		return runtime.AddError(fmt.Errorf("CANNOT_REMOVE_SELF"))
	}

	if err := runtime.Cursor(&models.ClientAdmin{}).Delete(&record).Error; err != nil {
		runtime.Debugf("destroy client admin error: %s", err.Error())
		return runtime.AddError(fmt.Errorf("SERVER_ERROR"))
	}

	runtime.Debugf("succesfully removed user[%d] as admin of [%d]", record.User, record.Client)

	return nil
}

func CreateClientAdmin(runtime *net.RequestRuntime) error {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.AddError(err)
	}

	validator := body.Validator()
	validator.Require("user")

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validator.HasErrors() == true {
		for _, m := range validator.Messages() {
			runtime.AddError(fmt.Errorf(m))
		}

		return nil
	}

	// attempt to parse out the user id from the body
	user, err := strconv.Atoi(body.Get("user"))

	if err != nil {
		return runtime.AddError(fmt.Errorf("INVALID_USER"))
	}

	// by default, we're only allowed to add users to the admin list of the current client
	client := runtime.Client.ID

	god := runtime.IsAdmin()

	// however, if the current user is a system admin, and a client has been provided, attempt to use it
	if god && body.KeyExists("client") {
		runtime.Debugf("admin attempting to make user %d admin of %v", user, body.Get("client"))
		input, err := strconv.Atoi(body.Get("client"))

		if err != nil {
			return runtime.AddError(fmt.Errorf("INVALID_CLIENT"))
		}

		client = uint(input)

		if err := runtime.Database().First(&models.Client{}, client).Error; err != nil {
			return runtime.AddError(fmt.Errorf("CLIENT_NOT_FOUND"))
		}
	}

	// if we are not a system admin, make sure we can even mess with the current client
	if god == false {
		admin := 0
		cursor := runtime.Cursor(&models.ClientAdmin{})
		if _ = cursor.Where("user = ? AND client = ?", runtime.User.ID, client).Count(&admin); admin == 0 {
			runtime.Debugf("unauthorized attempt to make user %d admin of %d", user, client)
			return runtime.AddError(fmt.Errorf("UNAUTHORIZED"))
		}
	}

	runtime.Debugf("attempting to add user %d as admin to client %d", user, client)
	mapping := models.ClientAdmin{User: uint(user), Client: client}

	dupe := 0
	cursor := runtime.Cursor(&models.ClientAdmin{})

	if _ = cursor.Where("user = ? AND client = ?", user, client).Count(&dupe); dupe != 0 {
		runtime.Debugf("duplicate entry: user %d with client %d", user, client)
		return runtime.AddError(fmt.Errorf("DUPLICATE_ENTRY"))
	}

	if err := runtime.Database().Create(&mapping).Error; err != nil {
		runtime.Debugf("unable to create: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED_SAVE"))
	}

	if client != runtime.Client.ID {
		manager := services.UserClientManager{runtime.Database()}

		u := models.User{Common: models.Common{ID: uint(user)}}
		c := models.Client{Common: models.Common{ID: client}}

		if _, err := manager.Associate(&u, &c); err != nil {
			runtime.Debugf("cant auto create admin clienttoken user[%d]-client[%d]: %s", user, client, err.Error())
		}
	}

	runtime.AddResult(mapping)

	return nil
}

func FindClientAdmins(runtime *net.RequestRuntime) error {
	var results []models.ClientAdmin
	blueprint := runtime.Blueprint()

	if runtime.IsAdmin() != true {
		runtime.Debugf("user is not admin, limiting query to client[%d]", runtime.Client.ID)
		err := blueprint.Filter("filter[client]", fmt.Sprintf("eq(%d)", runtime.Client.ID))

		if err != nil {
			runtime.Debugf("filter problem: %s", err.Error())
			return runtime.AddError(fmt.Errorf("PROBLEM"))
		}

		// make sure user is even able to see this client's admins by being a client admin themselces
		query := runtime.Database().Where("client = ? AND user = ?", runtime.Client.ID, runtime.User.ID)

		if err := query.Find(&results).Error; err != nil {
			runtime.Debugf("failed getting client admins for current situation problem: %s", err.Error())
			return runtime.AddError(fmt.Errorf("PROBLEM"))
		}

		if len(results) != 1 {
			runtime.Debugf("current user[%d] has no access to client[%d]", runtime.User.ID, runtime.Client.ID)
			return runtime.AddError(fmt.Errorf("NOT_FOUND"))
		}
	}

	total, err := blueprint.Apply(&results, runtime.Database())

	if err != nil {
		runtime.Debugf("BAD_LOOKUP: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_QUERY"))
	}

	for _, item := range results {
		runtime.AddResult(item)
	}

	runtime.SetMeta("total", total)

	return nil
}
