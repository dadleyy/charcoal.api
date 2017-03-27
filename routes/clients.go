package routes

import "fmt"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"
import "github.com/dadleyy/charcoal.api/util"

func FindClients(runtime *net.RequestRuntime) *net.ResponseBucket {
	blueprint := runtime.Blueprint()
	var clients []models.Client

	total, err := blueprint.Apply(&clients)

	if err != nil {
		runtime.Errorf("[find clients] unable to query clients: %s", err.Error())
		return runtime.ServerError()
	}

	meta := map[string]interface{}{"total": total}
	return &net.ResponseBucket{Results: clients, Meta: meta}
}

func UpdateClient(runtime *net.RequestRuntime) *net.ResponseBucket {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.SendErrors(fmt.Errorf("BAD_ID"))
	}

	if god := runtime.IsAdmin(); god != true {
		admin := models.ClientAdmin{}

		if e := runtime.Where("client_id = ? AND user_id = ?", id, runtime.User.ID).First(&admin).Error; e != nil {
			return runtime.LogicError("invalid-user")
		}
	}

	var client models.Client

	if err := runtime.First(&client, id).Error; err != nil {
		return runtime.LogicError("not-found")
	}

	body, err := forms.Parse(runtime.Request)

	if err != nil {
		runtime.Warnf("[update client] invalid body: %s", err.Error())
		return runtime.LogicError("invalid-body")
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
		return &net.ResponseBucket{Results: []interface{}{client}}
	}

	if err := runtime.Model(&client).Updates(updates).Error; err != nil {
		runtime.Errorf("[update client] failed updating client: %s", err.Error())
		return runtime.ServerError()
	}

	return &net.ResponseBucket{Results: []interface{}{client}}
}

func CreateClient(runtime *net.RequestRuntime) *net.ResponseBucket {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.SendErrors(fmt.Errorf("BAD_REQUEST"))
	}

	validator := body.Validator()

	validator.Require("name")
	validator.Require("description")
	validator.Require("redirect_uri")

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validator.HasErrors() == true {
		errors := []error{}

		for _, m := range validator.Messages() {
			errors = append(errors, fmt.Errorf("field:%s", m))
		}

		return runtime.SendErrors(errors...)
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
		runtime.Errorf("[create client] failed attempt to duplicate client: %v", err)
		return runtime.ServerError()
	}

	if err := cursor.Create(&client).Error; err != nil {
		runtime.Errorf("[create client] failed attempt to create client: %s", err.Error())
		return runtime.ServerError()
	}

	admin := models.ClientAdmin{ClientID: client.ID, UserID: runtime.User.ID}

	if err := runtime.Create(&admin).Error; err != nil {
		runtime.Errorf("[create client] failed automatically creating admin for client %d: %s", client.ID, err.Error())
		return runtime.SendErrors(err)
	}

	manager := services.UserClientManager{runtime.DB}

	if _, err := manager.Associate(&runtime.User, &client); err != nil {
		runtime.Errorf("[create client] err auto token user[%d]-client[%d]: %s", runtime.User.ID, client.ID, err.Error())
		return runtime.SendErrors(err)
	}

	return &net.ResponseBucket{Results: []interface{}{client}}
}
