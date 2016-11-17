package routes

import "fmt"
import "github.com/albrow/forms"
import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/services"

func FindClients(runtime *net.RequestRuntime) error {
	blueprint := runtime.Blueprint()
	var clients []models.Client

	total, err := blueprint.Apply(&clients, runtime.Database())

	if err != nil {
		runtime.Debugf("unable to query clients: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_QUERY"))
	}

	for _, client := range clients {
		runtime.AddResult(client)
	}

	runtime.SetMeta("total", total)

	return nil
}

func CreateClient(runtime *net.RequestRuntime) error {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_REQUEST"))
	}

	validator := body.Validator()

	validator.Require("name")
	validator.Require("description")
	validator.Require("redirect_uri")

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validator.HasErrors() == true {
		for _, m := range validator.Messages() {
			runtime.AddError(fmt.Errorf(m))
		}

		return nil
	}

	client := models.Client{
		Name:         body.Get("name"),
		Description:  body.Get("description"),
		RedirectUri:  body.Get("redirect_uri"),
		ClientID:     services.RandStringBytesMaskImprSrc(20),
		ClientSecret: services.RandStringBytesMaskImprSrc(40),
	}

	cursor := runtime.Database().Model(&client).Where("name = ?", client.Name)
	existing := 0

	if err := cursor.Count(&existing).Error; err != nil || existing >= 1 {
		runtime.Debugf("failed attempt to duplicate client: %v", err)
		return runtime.AddError(fmt.Errorf("INVALID_CLIENT_NAME"))
	}

	if err := cursor.Create(&client).Error; err != nil {
		runtime.Debugf("failed attempt to create client: %s", err.Error())
		return runtime.AddError(fmt.Errorf("SERVER_ERROR"))
	}

	admin := models.ClientAdmin{Client: client.ID, User: runtime.User.ID}

	if err := runtime.Database().Create(&admin).Error; err != nil {
		runtime.Debugf("failed automatically creating admin for client %d: %s", client.ID, err.Error())
		return runtime.AddError(err)
	}

	runtime.AddResult(client)
	return nil
}
