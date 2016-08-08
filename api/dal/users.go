package dal

import "fmt"
import "errors"
import "strings"
import "github.com/golang/glog"
import "golang.org/x/crypto/bcrypt"
import "github.com/asaskevich/govalidator"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/models"

type Updates map[string]interface{}

type UserFacade struct {
	*models.User
}

func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// validEmail
func validEmail(runtime *api.Runtime, target string, user models.User) bool {
	var existing models.User

	result := runtime.DB.Where("email = ?", target).Find(&existing)

	if missing := result.RecordNotFound(); !missing && user.ID != existing.ID {
		return false
	}

	return govalidator.IsEmail(target)
}


// FindUser
func FindUser(runtime *api.Runtime, blueprint *api.Blueprint) ([]models.User, int, error) {
	var users []models.User
	var total int

	head := blueprint.Apply(runtime)

	if e := head.Find(&users).Count(&total).Error; e != nil {
		glog.Errorf("error in FindUser: %s\n", e.Error())
		return users, -1, e
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

				// cast our hashed byte array to a string and update the stringval
				stringval = string(hashed)
			}

			if key == "email" && !validEmail(runtime, stringval, user) {
				glog.Errorf("attempted to update an email to invalud value on user %d\n", user.ID)
				return errors.New("invalid email")
			}

			// at this point we've validated the value as either a password, email or name. apply the
			// update to the database and continue onto the next key/value
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
func CreateUser(runtime *api.Runtime, facade *UserFacade) (models.User, error) {
	var user models.User

	if len(facade.Name) < 2 {
		return user, errors.New("invalid name")
	}

	if len(facade.Email) < 2 {
		return user, errors.New("invalid email")
	}

	if len(facade.Password) < 6 {
		return user, errors.New("passwords must be at least 6 characters long")
	}

	if valid := validEmail(runtime, facade.Email, user); !valid {
		return user, errors.New(fmt.Sprintf("invalid email: %s", facade.Email))
	}

	hashed, err := hash(facade.Password)

	if err != nil {
		return user, err
	}

	user = models.User{
		Email: facade.Email,
		Name: facade.Name,
		Password: string(hashed),
	}

	if err := runtime.DB.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
