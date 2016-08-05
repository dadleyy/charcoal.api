package user

import "errors"
import "github.com/meritoss/meritoss.api/api"
import "github.com/meritoss/meritoss.api/api/models"

func Create(runtime *api.Runtime, user *models.User) error {
	if len(user.Name) < 2 {
		return errors.New("invalid name")
	}

	if err := runtime.DB.Save(user).Error; err != nil {
		return err
	}

	return nil
}
