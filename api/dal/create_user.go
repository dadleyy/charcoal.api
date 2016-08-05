package dal

import "errors"
import "golang.org/x/crypto/bcrypt"

import "github.com/meritoss/meritoss.api/api"
import "github.com/meritoss/meritoss.api/api/models"

func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CreateUser(runtime *api.Runtime, user *models.User) error {
	if len(user.Name) < 2 {
		return errors.New("invalid name")
	}

	if len(user.Email) < 2 {
		return errors.New("invalid email")
	}

	if len(user.Password) < 6 {
		return errors.New("passwords must be at least 6 characters long")
	}

	hashed, err := hash(user.Password)

	if err != nil {
		return err
	}

	user.Password = string(hashed)

	if err := runtime.DB.Save(user).Error; err != nil {
		return err
	}

	return nil
}
