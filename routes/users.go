package routes

import "fmt"
import "github.com/labstack/echo"
import "golang.org/x/crypto/bcrypt"

import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/services"

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
		return fmt.Errorf("BAD_RUNTIME")
	}

	body, err := runtime.Body()

	if err != nil {
		runtime.Logger().Debug("unable to read request body for POST /users")
		return runtime.ErrorOut(err)
	}

	name, exists := body.String("name")

	if exists != true {
		return runtime.ErrorOut(fmt.Errorf("INVALID_NAME"))
	}

	email, exists := body.String("email")

	if exists != true {
		return runtime.ErrorOut(fmt.Errorf("INVALID_EMAIL"))
	}

	password, exists := body.String("password")

	if exists != true {
		return runtime.ErrorOut(fmt.Errorf("INVALID_PASSWORD"))
	}

	hashed, err := hash(password)

	if err != nil {
		return runtime.ErrorOut(fmt.Errorf("BAD_PASSWORD"))
	}

	user := models.User{
		Name: name,
		Email: email,
		Password: string(hashed),
	}

	if err = runtime.DB.Create(&user).Error; err != nil {
		return runtime.ErrorOut(err)
	}

	runtime.Logger().Debugf("successfully created user \"%s\", associating to client \"%d\"", name, runtime.Client.ID)

	manager := services.UserClientManager{runtime.DB}

	token, err := manager.AssociateClient(&user, &runtime.Client);

	if err != nil {
		runtime.Logger().Errorf("failed user[%d]-client[%d] associate: %s", user.ID, runtime.Client.ID, err.Error())
		return runtime.ErrorOut(fmt.Errorf("FAILED_ASSOCIATE"))
	}

	runtime.Logger().Debugf("created client/user token: %s", token.Token)
	runtime.Result(&user)

	return nil
}
