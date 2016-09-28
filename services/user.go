package services

import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

type UserManager struct {
	*context.Database
}

func (manager *UserManager) IsDuplicate(target *models.User) (bool, error) {
	var count int
	var existing models.User
	err := manager.Model(existing).Where(target).Count(&count).Error
	return count >= 1, err
}

func (manager *UserManager) FindOrCreate(target *models.User) error {
	return manager.FirstOrCreate(target, *target).Error
}
