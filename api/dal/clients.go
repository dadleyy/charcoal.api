package dal

import "errors"
import "strings"
import "crypto/rand"
import "encoding/hex"
import "github.com/golang/glog"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/db"
import "github.com/sizethree/meritoss.api/api/models"

type ClientFacade struct {
	Name string
}

// CreateClient
//
// given a database connection and a client facade, creates and returns a new client
// object, otherwise returns a non-nill error
func CreateClient(dbclient *db.Client, facade *ClientFacade) (models.Client, error) {
	var client models.Client

	name := strings.TrimSpace(facade.Name)

	if len(name) < 4 {
		return client, errors.New("client names must be at lease 4 characters long")
	}

	tokenbuffer, secretbuffer := make([]byte, 10), make([]byte, 20)

	if _, e := rand.Read(tokenbuffer); e != nil {
		return client, e
	}

	if _, e := rand.Read(secretbuffer); e != nil {
		return client, e
	}

	client = models.Client{
		Name: name,
		ClientSecret: hex.EncodeToString(secretbuffer),
		ClientID: hex.EncodeToString(tokenbuffer),
	}

	if e := dbclient.Save(&client).Error; e != nil {
		glog.Errorf("unable to save client: %s\n", e.Error())
		return client, e
	}

	return client, nil
}

// FindClient
// 
// given a database client and a blueprint, returns the list of appropriate clients
func FindClients(dbclient *db.Client, blueprint* api.Blueprint) ([]models.Client, int, error) {
	var clients []models.Client

	total, e := blueprint.Apply(&clients, dbclient)

	if e != nil {
		glog.Errorf("errror applying clients blueprint %s\n", e.Error())
		return clients, -1, e
	}

	return clients, total, nil
}

