package services

import "fmt"
import "errors"
import "encoding/hex"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

type UserClientManager struct {
	*context.Database
}

func (engine *UserClientManager) AssociateClient(user *models.User, client *models.Client) (models.ClientToken, error) {
	var result models.ClientToken

	if client.ID == 0 {
		return result, errors.New("BAD_CLIENT_ID")
	}

	buffer := []byte(fmt.Sprintf("%s:%s", client.ClientID, user.ID))

	result = models.ClientToken{
		Client: client.ID,
		User: user.ID,
	}

	var tcount uint

	if err := engine.Model(&result).Count(&tcount).Error; err != nil {
		return result, err
	}

	if tcount >= 1 {
		err := engine.Model(&result).Where(result).First(&result).Error
		return result, err
	}

	result.Token = hex.EncodeToString(buffer)

	if err := engine.Create(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}
