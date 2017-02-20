package routes

import "fmt"
import "strconv"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"

func CreateClientToken(runtime *net.RequestRuntime) error {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_REQUEST"))
	}

	// only allow system clients to "artificially" create client tokens
	if runtime.Client.System != true {
		return runtime.AddError(fmt.Errorf("BAD_CLIENT"))
	}

	validator := body.Validator()
	validator.Require("client")

	if validator.HasErrors() == true {
		for _, m := range validator.Messages() {
			runtime.AddError(fmt.Errorf(m))
		}

		return nil
	}

	target, err := strconv.Atoi(body.Get("client"))

	if err != nil {
		return runtime.AddError(fmt.Errorf("INVALID_CLIENT_ID"))
	}

	var client models.Client
	if err := runtime.Database().Where("id = ?", target).First(&client).Error; err != nil {
		return runtime.AddError(fmt.Errorf("CLIENT_NOT_FOUND"))
	}

	manager := services.UserClientManager{runtime.Database()}

	result, err := manager.Associate(&runtime.User, &client)

	if err != nil {
		runtime.Debugf("failed authorizing client %d for user %d: %s", client.ID, runtime.User.ID, err.Error())
		return runtime.AddError(fmt.Errorf("FAILED_ASSOCIATE"))
	}

	runtime.AddResult(result)
	return nil
}

func FindClientTokens(runtime *net.RequestRuntime) error {
	var tokens []models.ClientToken

	// start a db cursor based on the offset from the limit + page
	blueprint := runtime.Blueprint()

	// limit this query to to current user only
	total, err := blueprint.Apply(&tokens)

	if err != nil {
		runtime.Debugf("problem retreiving client tokens: %s", err.Error())
		return runtime.AddError(err)
	}

	for _, token := range tokens {
		runtime.AddResult(token)
	}

	runtime.SetMeta("total", total)

	return nil
}
