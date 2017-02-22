package routes

import "fmt"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"
import "github.com/dadleyy/charcoal.api/util"

func FindClients(runtime *net.RequestRuntime) error {
	blueprint := runtime.Blueprint()
	var clients []models.Client

	total, err := blueprint.Apply(&clients)

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

func UpdateClient(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.AddError(fmt.Errorf("BAD_ID"))
	}

	if god := runtime.IsAdmin(); god != true {
		admin := 0
		cursor := runtime.Cursor(&models.ClientAdmin{}).Where("client = ? AND user = ?", id, runtime.User.ID)

		if _ = cursor.Count(&admin); admin == 0 {
			return runtime.AddError(fmt.Errorf("UNAUTHORIZED: user[%d] client[%d]", runtime.User.ID, id))
		}
	}

	var client models.Client

	if err := runtime.First(&client, id).Error; err != nil {
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_REQUEST"))
	}

	updates := make(map[string]interface{})

	if body.KeyExists("description") {
		updates["description"] = body.Get("description")
	}

	if body.KeyExists("redirect_uri") {
		updates["redirect_uri"] = body.Get("redirect_uri")
	}

	if body.KeyExists("name") {
		updates["name"] = body.Get("name")
	}

	if len(updates) == 0 {
		runtime.AddResult(client)
		return nil
	}

	if err := runtime.Model(&client).Updates(updates).Error; err != nil {
		runtime.Debugf("failed updating client: %s", err.Error())
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	runtime.AddResult(client)

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
		ClientID:     util.RandStringBytesMaskImprSrc(20),
		ClientSecret: util.RandStringBytesMaskImprSrc(40),
	}

	cursor := runtime.Model(&client).Where("name = ?", client.Name)
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

	if err := runtime.Create(&admin).Error; err != nil {
		runtime.Debugf("failed automatically creating admin for client %d: %s", client.ID, err.Error())
		return runtime.AddError(err)
	}

	manager := services.UserClientManager{runtime.DB}

	if _, err := manager.Associate(&runtime.User, &client); err != nil {
		runtime.Debugf("failed auto token for user[%d]-client[%d]: %s", runtime.User.ID, client.ID, err.Error())
		return runtime.AddError(err)
	}

	runtime.AddResult(client)
	return nil
}
