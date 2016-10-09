package routes

import "fmt"
import "github.com/labstack/echo"
import "golang.org/x/crypto/bcrypt"

import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/services"

func hash(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(result), nil
}

func FindUser(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Runtime)

	blueprint := runtime.Blueprint()

	var users []models.User

	total, err := blueprint.Apply(&users, runtime.DB)

	if err != nil {
		runtime.Logger().Debugf("bad user lookup query: %s", err.Error())
		return runtime.ErrorOut(fmt.Errorf("BAD_QUERY"))
	}

	for _, user := range users {
		runtime.AddResult(&user)
	}

	runtime.AddMeta("total", total)

	return nil
}

func UpdateUser(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Runtime)
	id, err := runtime.ParamInt("id")

	if err != nil {
		runtime.Logger().Debugf("bad user id: %s", err.Error())
		return runtime.ErrorOut(fmt.Errorf("BAD_ID"))
	}


	if id != int(runtime.User.ID) {
		runtime.Logger().Debugf("invlaid user match request[%d]-runtime[%d]", id, runtime.User.ID)
		return runtime.ErrorOut(fmt.Errorf("BAD_ID"))
	}

	updates := make(map[string]interface{})
	applied := make(map[string]interface{})
	count := 0

	head := runtime.DB.Begin().Model(&runtime.User)

	if err := runtime.Bind(&updates); err != nil  {
		runtime.Logger().Debugf("bad update format: %s", err.Error())
		return runtime.ErrorOut(fmt.Errorf("BAD_FORMAT"))
	}

	if password, exists := updates["password"]; exists == true{
		password, ok := password.(string)

		if ok != true {
			return runtime.ErrorOut(fmt.Errorf("BAD_PASSWORD"))
		}

		if len(password) < 6 {
			return runtime.ErrorOut(fmt.Errorf("BAD_PASSWORD"))
		}

		password, err = hash(password)

		if err != nil {
			return runtime.ErrorOut(fmt.Errorf("BAD_PASSWORD"))
		}

		applied["password"] = password
		count++
	}

	if name, exists := updates["name"]; exists == true {
		name, ok := name.(string)

		if ok != true {
			return runtime.ErrorOut(fmt.Errorf("BAD_NAME"))
		}

		if len(name) < 2 {
			return runtime.ErrorOut(fmt.Errorf("BAD_NAME"))
		}

		applied["name"] = name
		count++
	}

	if email, exists := updates["email"]; exists == true {
		email, ok := email.(string)

		if ok != true {
			return runtime.ErrorOut(fmt.Errorf("BAD_EMAIL"))
		}

		if len(email) < 2 {
			return runtime.ErrorOut(fmt.Errorf("BAD_EMAIL"))
		}

		applied["email"] = email
		count++
	}

	if count == 0 {
		return nil
	}

	head.Updates(applied).Commit()

	runtime.Logger().Debugf("successfully updated user[%d]", id)
	runtime.AddResult(&runtime.User)
	return nil
}

func CreateUser(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Runtime)

	if ok != true {
		return fmt.Errorf("BAD_RUNTIME")
	}

	var target models.User

	if err := runtime.Bind(&target); err != nil  {
		runtime.Logger().Debugf("bad update format: %s", err.Error())
		return runtime.ErrorOut(fmt.Errorf("BAD_FORMAT"))
	}

	if target.Name == nil || len(*target.Name) < 2 {
		return runtime.ErrorOut(fmt.Errorf("BAD_NAME"))
	}

	if target.Email == nil || len(*target.Email) < 2 {
		return runtime.ErrorOut(fmt.Errorf("BAD_EMAIL"))
	}

	if target.Password == nil || len(*target.Password) < 6 {
		return runtime.ErrorOut(fmt.Errorf("BAD_PASSWORD"))
	}

	usermgr := services.UserManager{runtime.DB}

	if dupe, err := usermgr.IsDuplicate(&target); dupe || err != nil {
		runtime.Logger().Debugf("duplicate user")
		return runtime.ErrorOut(fmt.Errorf("BAD_USER"))
	}

	hashed, err := hash(*target.Password)

	if err != nil {
		return runtime.ErrorOut(fmt.Errorf("BAD_PASSWORD"))
	}

	target.Password = &hashed

	if err := runtime.DB.Create(&target).Error; err != nil {
		runtime.Logger().Debugf("unable to save: %s", err.Error())
		return runtime.ErrorOut(fmt.Errorf("FAILED"))
	}

	clientmgr := services.UserClientManager{runtime.DB}

	token, err := clientmgr.Associate(&target, &runtime.Client)

	if err != nil {
		runtime.Logger().Debugf("unable to associate: %s", err.Error())
		return runtime.ErrorOut(fmt.Errorf("FAILED"))
	}

	runtime.Logger().Debugf("associated user[%d] with client[%d]: %s", target.ID, runtime.Client.ID, token.Token)
	runtime.AddResult(&target)

	return nil
}
