package dal

import "fmt"
import "errors"
import "strings"
import "crypto/rand"
import "encoding/hex"
import "github.com/golang/glog"
import "golang.org/x/crypto/bcrypt"
import "github.com/asaskevich/govalidator"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/db"
import "github.com/sizethree/meritoss.api/api/models"

type UserFacade struct {
	Name string
	Email string
	Password string
}

func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// validEmail
func validEmail(client *db.Client, target string, user models.User) bool {
	var existing models.User

	result := client.Where("email = ?", target).Find(&existing)

	if missing := result.RecordNotFound(); !missing && user.ID != existing.ID {
		return false
	}

	return govalidator.IsEmail(target)
}

// FindUser
// 
// given a database client and a blueprint, returns an array of models, an integer representing
// the total count of users matched by the blueprint and optionally an error
func FindUser(client *db.Client, blueprint *api.Blueprint) ([]models.User, int, error) {
	var users []models.User

	total, e := blueprint.Apply(&users, client)

	if e != nil {
		glog.Errorf("error in FindUser: %s\n", e.Error())
		return users, -1, e
	}

	return users, total, nil
}

// UpdateUser
//
// given a database client, this function attempts to load in
func UpdateUser(client *db.Client, updates *Updates, userid int) error {
	var user models.User

	head := client.Where("ID = ?", userid).Find(&user)

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

			if key == "email" && !validEmail(client, stringval, user) {
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

func AuthorizeClient(dbclient *db.Client, userid uint, clientid uint) error {
	var client models.Client
	var user models.User

	if result := dbclient.Where("ID = ?", userid).First(&user); result.RecordNotFound() || result.Error != nil {
		return result.Error
	}

	if result := dbclient.Where("ID = ?", clientid).First(&client); result.RecordNotFound() || result.Error != nil {
		return result.Error
	}

	// generate a random token buffer of 20 character length (10 bytes * 2 hex characters per byte)
	tokenbuffer := make([]byte, 10)
	_, err := rand.Read(tokenbuffer)

	if err != nil {
		return err
	}

	newtoken := models.ClientToken{
		Client: client.ID,
		User: user.ID,
		Token: hex.EncodeToString(tokenbuffer),
	}

	if e := dbclient.Save(&newtoken).Error; e != nil {
		return e
	}

	glog.Infof("found user %d and client %d, generating token: %s\n", user.ID, client.ID, newtoken.Token)

	return nil
}

// CreateUser
func CreateUser(client *db.Client, facade *UserFacade) (models.User, error) {
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

	if valid := validEmail(client, facade.Email, user); !valid {
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

	if err := client.Save(&user).Error; err != nil {
		return user, err
	}

	if err = AuthorizeClient(client, user.ID, 1); err != nil {
		return user, err
	}

	return user, nil
}
