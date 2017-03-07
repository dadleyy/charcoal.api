package services

import "fmt"
import "regexp"
import "net/url"
import "strings"
import "unicode"
import "net/mail"
import "github.com/jinzhu/gorm"
import "golang.org/x/crypto/bcrypt"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/models"

const UserManagerErrorUnauthorizedDomain = "unauthorized-domain"
const UserManagerErrorDuplicate = "duplicate-user"
const UserManagerUsernameRE = "^[a-zA-Z0-9\\-_]{6,20}$"

type UserManager struct {
	*gorm.DB
	*log.Logger
}

func (manager *UserManager) ValidUsername(username string) bool {
	re := regexp.MustCompile(UserManagerUsernameRE)
	return re.MatchString(username)
}

func (manager *UserManager) ValidPassword(password string) bool {
	if len(password) < 6 || len(password) > 20 {
		return false
	}

	for _, c := range password {
		if unicode.IsSpace(c) {
			return false
		}
	}

	return true
}

func (manager *UserManager) ValidUser(user *models.User) (bool, []error) {
	errors := make([]error, 0)

	if manager.ValidDomain(user.Email) != true {
		errors = append(errors, fmt.Errorf("reason:%s", UserManagerErrorUnauthorizedDomain))
	}

	if dupe, err := manager.IsDuplicate(user); dupe || err != nil {
		errors = append(errors, fmt.Errorf("reason:%s", UserManagerErrorDuplicate))
	}

	return len(errors) == 0, errors
}

func (manager *UserManager) IsDuplicate(target *models.User) (bool, error) {
	var count int
	var existing models.User

	if target == nil {
		return true, fmt.Errorf("INVALID_USER")
	}

	err := manager.Model(existing).Where("email = ? OR username = ?", target.Email, target.Username).Count(&count).Error
	return count >= 1, err
}

func (manager *UserManager) FindOrCreate(target *models.User) error {
	if target == nil {
		return fmt.Errorf("BAD_TARGET")
	}

	return manager.Where(models.User{Email: target.Email}).FirstOrCreate(target).Error
}

func (manager *UserManager) ApplyUpdates(existing *models.User, updates url.Values) []error {
	errors := make([]error, 0)
	applied := make(map[string]interface{})

	if _, ok := updates["password"]; ok == true {
		password := updates.Get("password")
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			errors = append(errors, err)
		} else if manager.ValidPassword(password) != true {
			manager.Debugf("password [%s] does not pass validation", password)
			errors = append(errors, fmt.Errorf("field:password"))
		} else {
			manager.Debugf("updating password on user[%d]", existing.ID)
			applied["password"] = string(hashed)
		}
	}

	if _, ok := updates["name"]; ok == true && updates.Get("name") != existing.Name {
		applied["name"] = updates.Get("name")
	}

	if _, ok := updates["email"]; ok == true && updates.Get("email") != existing.Email {
		desired := updates.Get("email")

		if dupe, err := manager.IsDuplicate(&models.User{Email: desired}); err != nil || dupe {
			manager.Debugf("email[%s] considered duplicate", desired)
			errors = append(errors, fmt.Errorf("reason:invalid-email"))
		} else if _, err := mail.ParseAddress(desired); err != nil {
			manager.Debugf("email[%s] did not pass inspection: %s", desired, err.Error())
			errors = append(errors, fmt.Errorf("reason:invalid-email"))
		} else {
			applied["email"] = updates.Get("email")
		}
	}

	if _, ok := updates["username"]; ok == true && updates.Get("username") != existing.Username {
		desired := updates.Get("username")

		if dupe, err := manager.IsDuplicate(&models.User{Username: desired}); err != nil || dupe {
			manager.Debugf("username[%s] considered a duplicate", desired)
			errors = append(errors, fmt.Errorf("reason:invalid-username"))
		} else if manager.ValidUsername(desired) != true {
			manager.Debugf("username[%s] did not pass inspection", desired)
			errors = append(errors, fmt.Errorf("reason:invalid-username"))
		} else {
			applied["username"] = updates.Get("username")
		}
	}

	if len(errors) >= 1 {
		return errors
	}

	if e := manager.Model(existing).Updates(applied).Error; e != nil {
		manager.Warnf("failed update to user[%d]: %s", existing.ID, e.Error())
		return []error{fmt.Errorf("reason:server-error")}
	}

	return errors
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
