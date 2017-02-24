package services

import "fmt"
import "strings"
import "github.com/jinzhu/gorm"
import "github.com/dadleyy/charcoal.api/models"

const ErrUnauthorizedDomain = "UNAUTHORIZED_DOMAIN"

type UserManager struct {
	*gorm.DB
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

func (manager *UserManager) ValidDomain(email string) bool {
	email = strings.TrimSpace(email)
	var settings models.SystemSettings

	// no settings found - we're good to move on
	if err := manager.First(&settings).Error; err != nil {
		return true
	}

	if settings.RestrictedEmailDomains == false {
		return true
	}

	var whitelist []models.SystemEmailDomain

	if err := manager.Find(&whitelist).Error; err != nil {
		return false
	}

	last := strings.LastIndex(email, "@")

	if last == -1 {
		return false
	}

	if last+1 >= len(email) {
		return false
	}

	domain := email[last+1:]

	for _, allowed := range whitelist {
		if allowed.Domain == domain {
			return true
		}
	}

	return false
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
