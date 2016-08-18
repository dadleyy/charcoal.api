package dal

import "errors"
import "crypto/rand"
import "encoding/hex"

import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/models"

type ClientTokenFacade struct {
	Client uint
	User uint
	Referrer models.Client
}

// duplicateToken
// 
// Helper function to return the existence of a client token with the same user id and client
// id as those sent into the function
func duplicateToken(dbclient *db.Client, user uint, client uint) bool {
	var existing models.ClientToken
	result := dbclient.Where("user = ? AND client = ?", user, client).First(&existing)
	return result.RecordNotFound() != true
}

// RandomToken
// 
// Helper function to generate a random hex token of `size` length (in bytes).
func RandomToken(size int) (string, error) {
	tokenbuffer := make([]byte, size)

	if _, err := rand.Read(tokenbuffer); err != nil {
		return "", err
	}

	return hex.EncodeToString(tokenbuffer), nil
}

func CreateClientToken(dbclient *db.Client, facade *ClientTokenFacade) (models.ClientToken, error) {
	var token models.ClientToken
	var client models.Client

	if facade.User < 1 {
		return token, errors.New("invalid user")
	}

	// will eventually want to support "referrer" clients that can behave as the main miritos client
	// but for now this is fine.
	if facade.Referrer.ID != 1 {
		return token, errors.New("unauthorized client")
	}

	// check to make sure we are dealing with a valid client
	if e := dbclient.Where("id = ?", facade.Client).First(&client).Error; e != nil {
		return token, e
	}

	// check to make sure this combination of user/client has not already had a token created for it
	if duplicateToken(dbclient, facade.User, client.ID) {
		return token, errors.New("duplicate token for user/client")
	}

	tokenstr, e := RandomToken(10)

	if e != nil {
		return token, e
	}

	token = models.ClientToken{Token: tokenstr, User: facade.User, Client: client.ID}

	if e := dbclient.Save(&token).Error; e != nil {
		return token, e
	}

	return token, nil
}

func ClientTokensForUser(dbclient *db.Client, user uint) ([]models.ClientToken, int, error) {
	var tokens []models.ClientToken
	var total int

	result := dbclient.Where("user = ?", user).Find(&tokens).Count(&total)

	if e := result.Error; e != nil {
		return tokens, total, e
	}

	return tokens, total, nil
}
