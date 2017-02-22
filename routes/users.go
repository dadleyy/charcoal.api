package routes

import "fmt"
import "github.com/albrow/forms"
import "golang.org/x/crypto/bcrypt"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"

func hash(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(result), nil
}

func cleanseUser(user models.User) interface{} {
	return struct {
		models.Common
		Name  string `json:"name"`
		Email string `json:"email"`
	}{user.Common, *user.Name, *user.Email}
}

func CreateUser(runtime *net.RequestRuntime) error {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.AddError(err)
	}

	validator := body.Validator()

	validator.Require("email")
	validator.MatchEmail("email")

	validator.Require("password")
	validator.LengthRange("password", 6, 20)

	validator.Require("name")
	validator.LengthRange("password", 2, 100)

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validator.HasErrors() == true {
		for _, m := range validator.Messages() {
			runtime.AddError(fmt.Errorf(m))
		}

		return nil
	}

	password, err := hash(body.Get("password"))

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_USER"))
	}

	email := body.Get("email")
	name := body.Get("name")

	user := models.User{Email: &email, Password: &password, Name: &name}

	usrmgr := services.UserManager{runtime.DB}

	if usrmgr.ValidDomain(email) != true {
		runtime.Debugf("attempt to sign up w/ invalid domain: %s", email)
		return runtime.AddError(fmt.Errorf(services.ErrUnauthorizedDomain))
	}

	if dupe, err := usrmgr.IsDuplicate(&user); dupe || err != nil {
		runtime.Debugf("duplicate user")
		return runtime.AddError(fmt.Errorf("BAD_USER"))
	}

	if err := runtime.Create(&user).Error; err != nil {
		runtime.Debugf("unable to save: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED"))
	}

	clientmgr := services.UserClientManager{runtime.DB}

	if _, err := clientmgr.Associate(&user, &runtime.Client); err != nil {
		runtime.Debugf("unable to associate: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED"))
	}

	runtime.Debugf("associated user[%d] with client[%d]", user.ID, runtime.Client.ID)
	runtime.AddResult(cleanseUser(user))

	return nil
}

func UpdateUser(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		runtime.AddError(fmt.Errorf("BAD_ID"))
		return nil
	}

	if runtime.User.ID != uint(id) {
		runtime.AddError(fmt.Errorf("BAD_ID"))
		return nil
	}

	body, err := forms.Parse(runtime.Request)

	if err != nil {
		runtime.AddError(err)
		return nil
	}

	validate := body.Validator()
	updates := make(map[string]interface{})

	// if an email is present, validate it
	if body.KeyExists("email") {
		validate.Require("email")
		validate.MatchEmail("email")

		email := body.Get("email")
		current := *(runtime.User.Email)

		updates["email"] = email

		manager := services.UserManager{runtime.DB}

		if dupe, err := manager.IsDuplicate(&models.User{Email: &email}); (email != current) && (err != nil || dupe) {
			return runtime.AddError(fmt.Errorf("BAD_EMAIL"))
		}
	}

	// if a password is present, validate it
	if body.KeyExists("password") {
		validate.Require("password")
		validate.LengthRange("password", 6, 20)
		password := body.Get("password")

		hashed, err := hash(password)

		if err != nil {
			runtime.AddError(fmt.Errorf("BAD_PASSWORD"))
			return nil
		}

		updates["password"] = hashed
	}

	// if a password is present, validate it
	if body.KeyExists("name") {
		validate.Require("name")
		validate.LengthRange("name", 2, 100)
		updates["name"] = body.Get("name")
	}

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validate.HasErrors() == true {
		for _, m := range validate.Messages() {
			runtime.AddError(fmt.Errorf(m))
		}

		return nil
	}

	if err := runtime.Model(&runtime.User).Updates(updates).Error; err != nil {
		runtime.AddError(err)
		return nil
	}

	runtime.AddResult(cleanseUser(runtime.User))
	runtime.Debugf("updating user[%d]", id)
	return nil
}

func FindUser(runtime *net.RequestRuntime) error {
	blue := runtime.Blueprint()
	var users []models.User

	count, err := blue.Apply(&users)

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_BLUEPRINT"))
	}

	runtime.SetMeta("count", count)

	for _, u := range users {
		runtime.AddResult(cleanseUser(u))
	}

	return nil
}
