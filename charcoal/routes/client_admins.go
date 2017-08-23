package routes

import "fmt"

import "strconv"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/charcoal/net"
import "github.com/dadleyy/charcoal.api/charcoal/models"
import "github.com/dadleyy/charcoal.api/charcoal/services"

func DeleteClientAdmin(runtime *net.RequestRuntime) *net.ResponseBucket {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.SendErrors(fmt.Errorf("BAD_ID"))
	}

	var record models.ClientAdmin

	if err := runtime.First(&record, id).Error; err != nil {
		return runtime.SendErrors(fmt.Errorf("NOT_FOUND"))
	}

	authorized := runtime.IsAdmin()

	// if the user is not a system admin, check to see if they are an admin of the client
	if authorized != true {
		count := 0
		cursor := runtime.Model(&models.ClientAdmin{}).Where("client = ? AND user = ?", record.ClientID, runtime.User.ID)

		if err := cursor.Count(&count).Error; err != nil || count == 0 {
			message := "[del client admin] unauthorized attempt to remove client admin user[%d] record[%d]: %v"
			runtime.Warnf(message, runtime.User.ID, id, err)
			return runtime.SendErrors(fmt.Errorf("UNAUTHORIZED"))
		}
	}

	if record.UserID == runtime.User.ID {
		return runtime.SendErrors(fmt.Errorf("CANNOT_REMOVE_SELF"))
	}

	if err := runtime.Model(&models.ClientAdmin{}).Delete(&record).Error; err != nil {
		runtime.Errorf("[del client admin] destroy client admin error: %s", err.Error())
		return runtime.SendErrors(fmt.Errorf("SERVER_ERROR"))
	}

	runtime.Debugf("[del client admin] succesfully removed user[%d] as admin of [%d]", record.UserID, record.ClientID)

	return &net.ResponseBucket{}
}

func CreateClientAdmin(runtime *net.RequestRuntime) *net.ResponseBucket {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.SendErrors(err)
	}

	validator := body.Validator()
	validator.Require("user_id")

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validator.HasErrors() == true {
		errors := []error{}

		for _, m := range validator.Messages() {
			errors = append(errors, fmt.Errorf("field:%s", m))
		}

		return runtime.SendErrors(errors...)
	}

	// attempt to parse out the user id from the body
	user, err := strconv.Atoi(body.Get("user_id"))

	if err != nil {
		return runtime.SendErrors(fmt.Errorf("INVALID_USER"))
	}

	// by default, we're only allowed to add users to the admin list of the current client
	client := runtime.Client.ID

	god := runtime.IsAdmin()

	// however, if the current user is a system admin, and a client has been provided, attempt to use it
	if god && body.KeyExists("client_id") {
		runtime.Debugf("[create client admin] admin attempting to make user %d admin of %v", user, body.Get("client"))
		input, err := strconv.Atoi(body.Get("client_id"))

		if err != nil {
			return runtime.LogicError("invalid-client")
		}

		client = uint(input)

		if err := runtime.First(&models.Client{}, client).Error; err != nil {
			return runtime.LogicError("missing-client")
		}
	}

	// if we are not a system admin, make sure we can even mess with the current client
	if god == false {
		admin := 0
		cursor := runtime.Model(&models.ClientAdmin{})

		if _ = cursor.Where("user_id = ? AND cient_id = ?", runtime.User.ID, client).Count(&admin); admin == 0 {
			runtime.Warnf("[create client admin] unauthorized attempt to make user %d admin of %d", user, client)
			return runtime.LogicError("not-found")
		}
	}

	runtime.Debugf("[create client admin] attempting to add user %d as admin to client %d", user, client)
	mapping := models.ClientAdmin{UserID: uint(user), ClientID: client}

	dupe := 0
	cursor := runtime.Model(&models.ClientAdmin{})

	if _ = cursor.Where("user_id = ? AND client_id = ?", user, client).Count(&dupe); dupe != 0 {
		runtime.Warnf("[create client admin] duplicate entry: user %d with client %d", user, client)
		return runtime.LogicError("duplicate-entry")
	}

	if err := runtime.Create(&mapping).Error; err != nil {
		runtime.Errorf("[create client admin] unable to create: %s", err.Error())
		return runtime.ServerError()
	}

	if client != runtime.Client.ID {
		manager := services.UserClientManager{runtime.DB}

		u := models.User{Common: models.Common{ID: uint(user)}}
		c := models.Client{Common: models.Common{ID: client}}

		if _, err := manager.Associate(&u, &c); err != nil {
			runtime.Debugf("[create client admin] failed auto token user[%d]-client[%d]: %s", user, client, err.Error())
		}
	}

	return &net.ResponseBucket{Results: []interface{}{mapping}}
}

func FindClientAdmins(runtime *net.RequestRuntime) *net.ResponseBucket {
	var results []models.ClientAdmin
	blueprint := runtime.Blueprint()

	// if the runtime is not operating under admin privileges
	if runtime.IsAdmin() != true {
		runtime.Debugf("[find client admin] user is not admin, limiting query to client[%d]", runtime.Client.ID)

		blueprint = runtime.Blueprint(runtime.Where("client_id = ?", runtime.Client.ID))

		// make sure user is even able to see this client's admins by being a client admin themselces
		query := runtime.Where("client_id = ? AND user_id = ?", runtime.Client.ID, runtime.User.ID)

		if err := query.Find(&results).Error; err != nil {
			runtime.Errorf("[find client admin] failed getting client admins for current situation problem: %s", err.Error())
			return runtime.ServerError()
		}

		if len(results) != 1 {
			runtime.Debugf("[fiend client admin] bad user[%d]: no access client[%d]", runtime.User.ID, runtime.Client.ID)
			return runtime.ServerError()
		}
	}

	total, err := blueprint.Apply(&results)

	if err != nil {
		runtime.Errorf("[find admins] bad blueprint: %s", err.Error())
		return runtime.ServerError()
	}

	meta := map[string]interface{}{"total": total}
	return &net.ResponseBucket{Results: []interface{}{results}, Meta: meta}
}
