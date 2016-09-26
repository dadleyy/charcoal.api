package routes

import "errors"
import "github.com/labstack/echo"
import "golang.org/x/crypto/bcrypt"

import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func FindUser(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Miritos)
	blueprint := runtime.Blueprint()
	var users []models.User

	total, err := blueprint.Apply(&users, runtime.DB)

	if err != nil {
		return err
	}


	for _, user := range users {
		runtime.Result(&user)
	}

	runtime.SetMeta("total", total)

	return nil
}

func UpdateUser(ectx echo.Context) error {
	return nil
}

func CreateUser(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Miritos)

	if ok != true {
		return errors.New("unable to load miritos context")
	}

	body, err := runtime.Body()

	if err != nil {
		runtime.Error(err)
		return err
	}

	name, exists := body.String("name")

	if exists != true {
		runtime.Error(errors.New("must provide a valid \"name\""))
		return nil
	}

	email, exists := body.String("email")

	if exists != true {
		runtime.Error(errors.New("must provide a valid \"email\""))
		return nil
	}

	password, exists := body.String("email")

	if exists != true {
		runtime.Error(errors.New("must provide a valid \"password\""))
		return nil
	}

	hashed, err := hash(password)

	if err != nil {
		runtime.Error(errors.New("unable to hash password"))
		return nil
	}

	user := models.User{
		Name: name,
		Email: email,
		Password: string(hashed),
	}

	if err = runtime.DB.Create(&user).Error; err != nil {
		runtime.Error(err)
		return nil
	}

	runtime.Result(&user)
	runtime.Logger().Infof("attempting to create user \"%s\", pass: %s", name, password)

	return nil
}
