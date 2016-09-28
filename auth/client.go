package auth

import "fmt"
import "errors"
import "encoding/hex"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

type Client struct {
	*context.Database
}

func (engine *Client) IsDuplicate(user *models.User) (bool, error) {
	var count int
	var existing models.User
	err := engine.Model(&existing).Where("email = ?", user.Email).Count(&count).Error
	return count >= 1, err
}

func (engine *Client) AssociateClient(user *models.User, client *models.Client) error {
	dupe, err := engine.IsDuplicate(user)

	if err != nil {
		return errors.New("BAD_USER")
	}

	if dupe != true {
		err = engine.Create(user).Error
	}

	if err != nil {
		return err
	}

	err = engine.Where("email = ?", user.Email).First(user).Error

	if err != nil {
		return err
	}

	if client.ID == 0 {
		return errors.New("BAD_CLIENT_ID")
	}

	buffer := []byte(fmt.Sprintf("%s:%s", client.ClientID, user.ID))

	token := models.ClientToken{
		Client: client.ID,
		User: user.ID,
	}

	var tcount uint

	if err := engine.Model(&token).Count(&tcount).Error; err != nil {
		return err
	}

	if tcount >= 1 {
		return nil
	}

	token.Token = hex.EncodeToString(buffer)

	if err := engine.Create(&token).Error; err != nil {
		return err
	}

	return nil
}
