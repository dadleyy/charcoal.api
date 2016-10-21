package services

import "fmt"
import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/models"

type UserManager struct {
	*db.Connection
}

func (manager *UserManager) IsDuplicate(target *models.User) (bool, error) {
	var count int
	var existing models.User

	if target == nil {
		return true, fmt.Errorf("INVALID_USER")
	}

	err := manager.Model(existing).Where("email = ?", target.Email).Count(&count).Error
	return count >= 1, err
}

func (manager *UserManager) FindOrCreate(target *models.User) error {
	if target == nil {
		return fmt.Errorf("BAD_TARGET")
	}

	return manager.Where(models.User{Email: target.Email}).FirstOrCreate(target).Error
}

func (manager *UserManager) IsAdmin(target *models.User) bool {
	if target == nil || target.ID == 0 {
		return false
	}

	var maps []models.UserRoleMapping

	if err := manager.Where("user = ?", target.ID).Find(&maps).Error; err != nil {
		return false
	}

	for _, mapping := range maps {
		if mapping.Role == 1 {
			return true
		}
	}

	return false
}
