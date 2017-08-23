package routes

import "fmt"
import "strconv"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/charcoal/net"
import "github.com/dadleyy/charcoal.api/charcoal/models"
import "github.com/dadleyy/charcoal.api/charcoal/services"

func CreateClientToken(runtime *net.RequestRuntime) *net.ResponseBucket {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.SendErrors(fmt.Errorf("BAD_REQUEST"))
	}

	// only allow system clients to "artificially" create client tokens
	if runtime.Client.System != true {
		return runtime.SendErrors(fmt.Errorf("BAD_CLIENT"))
	}

	validator := body.Validator()
	validator.Require("client")

	if validator.HasErrors() == true {
		errors := []error{}

		for _, m := range validator.Messages() {
			errors = append(errors, fmt.Errorf("field:%s", m))
		}

		return runtime.SendErrors(errors...)
	}

	target, err := strconv.Atoi(body.Get("client"))

	if err != nil {
		return runtime.SendErrors(fmt.Errorf("INVALID_CLIENT_ID"))
	}

	var client models.Client
	if err := runtime.Where("id = ?", target).First(&client).Error; err != nil {
		return runtime.SendErrors(fmt.Errorf("CLIENT_NOT_FOUND"))
	}

	manager := services.UserClientManager{runtime.DB}

	result, err := manager.Associate(&runtime.User, &client)

	if err != nil {
		runtime.Errorf("[create token] client %d for user %d: %s", client.ID, runtime.User.ID, err.Error())
		return runtime.SendErrors(fmt.Errorf("FAILED_ASSOCIATE"))
	}

	return &net.ResponseBucket{Results: []interface{}{result}}
}

func FindClientTokens(runtime *net.RequestRuntime) *net.ResponseBucket {
	var tokens []models.ClientToken

	// start a db cursor based on the offset from the limit + page
	blueprint := runtime.Blueprint()

	// limit this query to to current user only
	total, err := blueprint.Apply(&tokens)

	if err != nil {
		runtime.Debugf("[find tokens] problem retreiving client tokens: %s", err.Error())
		return runtime.SendErrors(err)
	}

	meta := map[string]interface{}{"total": total}
	return &net.ResponseBucket{Results: tokens, Meta: meta}
}
