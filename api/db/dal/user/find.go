package user

import "github.com/meritoss/meritoss.api/api"
import "github.com/meritoss/meritoss.api/api/models"

func Find(runtime api.Runtime) ([]models.User, error) {
	var users []models.User;

	runtime.DB.Where("id > 0").Where("name = ? ", "danny").Find(&users)

	return users, nil
}
