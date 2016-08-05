package dal

import "fmt"
import "errors"
import "strings"
import "github.com/golang/glog"
import "golang.org/x/crypto/bcrypt"
import "github.com/asaskevich/govalidator"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/models"
import "github.com/sizethree/meritoss.api/api/middleware"

type Updates map[string]interface{}

func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// validEmail
func validEmail(runtime *api.Runtime, user models.User) bool {
	var existing models.User

	result := runtime.DB.Where("email = ?", user.Email).Find(&existing)

	if missing := result.RecordNotFound(); !missing && user.ID != existing.ID {
		return false
	}

	return govalidator.IsEmail(user.Email)
}


// FindUser
func FindUser(runtime api.Runtime, blueprint middleware.Blueprint) ([]models.User, int, error) {
	var users []models.User
	var total int

	limit, offset := blueprint.Limit, blueprint.Limit * blueprint.Page

	head := runtime.DB.Begin().Limit(limit).Offset(offset)

	for _, filter := range blueprint.Filters {
		glog.Infof("adding filter \"%s\"\n", filter.Reduce())
		head = head.Where(filter.Reduce())
	}

	result := head.Find(&users).Count(&total).Commit()

	if result.Error != nil {
		glog.Errorf("error in FindUser: %s\n", result.Error.Error())
		return users, -1, result.Error
	}

	return users, total, nil
}

// UpdateUser
func UpdateUser(runtime *api.Runtime, updates *Updates, userid int) error {
	var user models.User

	head := runtime.DB.Where("ID = ?", userid).Find(&user)

	for key, value := range *updates {
		key = strings.ToLower(strings.TrimSpace(key))

		// if we're receiving a password in our data, we need to convert to string and hash
		if key == "password" || key == "email" || key == "name" {
			stringval, ok := value.(string)

			if !ok {
				return errors.New(fmt.Sprintf("bad value for %s", key))
			}

			if key == "password" {
				if len(stringval) < 6 {
					return errors.New("password must be at least 6 characters long")
				}

				hashed, err := hash(stringval)

				if err != nil {
					return errors.New("unable to hash password")
				}

				head = head.Update(key, hashed)
				continue
			}

			if key == "email" && !validEmail(runtime, user) {
				return errors.New("invalid email")
			}

			head = head.Update(key, stringval)
		}
	}

	if e := head.Error; e != nil {
		glog.Errorf("unable to perform update: %s\n", e.Error())
		return e
	}

	return nil
}


// CreateUser
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

	if valid := validEmail(runtime, *user); !valid {
		return errors.New(fmt.Sprintf("invalid email: %s", user.Email))
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
