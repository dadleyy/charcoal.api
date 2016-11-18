package routes

import "fmt"

import "strconv"
import "github.com/albrow/forms"
import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"

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
		runtime.Debugf("admin attempting to make user %d admin of %v", body.Get("user"), body.Get("client"))
		input, err := strconv.Atoi(body.Get("client"))

		if err != nil {
			return runtime.AddError(fmt.Errorf("INVALID_CLIENT"))
		}

		client = uint(input)
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

	runtime.Debugf("attempting to add user %d as admin to client %d", user, runtime.Client.ID)
	mapping := models.ClientAdmin{User: uint(user), Client: runtime.Client.ID}

	dupe := 0
	cursor := runtime.Cursor(&models.ClientAdmin{})

	if _ = cursor.Where("user = ? AND client = ?", user, runtime.Client.ID).Count(&dupe); dupe != 0 {
		runtime.Debugf("duplicate entry: user %d with client %d", user, runtime.Client.ID)
		return runtime.AddError(fmt.Errorf("DUPLICATE_ENTRY"))
	}

	if err := runtime.Database().Create(&mapping).Error; err != nil {
		runtime.Debugf("unable to create: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED_SAVE"))
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
