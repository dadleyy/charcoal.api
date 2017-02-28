package routes

import "fmt"
import "bytes"
import "testing"
import "github.com/jinzhu/gorm"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/routes/testutils"

func clean(db *gorm.DB) {
	db.Exec("DELETE FROM user_role_mappings where id > 1")
	db.Exec("DELETE FROM client_admins where id > 1")
	db.Exec("DELETE FROM clients where id > 1")
	db.Exec("DELETE FROM users where id > 1")
	db.Exec("DELETE FROM client_tokens where id > 1")
}

func Test_UpdateClient_GodUser(t *testing.T) {
	reader := bytes.NewReader([]byte("{\"name\": \"updated-name\"}"))

	client, user := models.Common{ID: 1002}, models.Common{ID: 2002}
	context := testutils.New("PATCH", "clients/:id", fmt.Sprintf("/clients/%d", client.ID), "application/json", reader)

	defer context.Database.Close()
	defer clean(context.Database)

	info := struct {
		name  string
		email string
	}{"test1", "test@test.com"}

	context.Request.Client = models.Client{
		Common:       client,
		Name:         "clients-test1",
		ClientID:     "123123123",
		ClientSecret: "client-admin-routes-secret",
		System:       true,
	}

	context.Request.User = models.User{
		Common: user,
		Name:   &info.name,
		Email:  &info.email,
	}

	context.Database.Create(&context.Request.Client)
	context.Database.Create(&context.Request.User)
	context.Database.Create(&models.UserRoleMapping{User: context.Request.User.ID, Role: 1})

	err := UpdateClient(&context.Request)

	if err != nil {
		t.Fatalf("User w/ God privileges was unable to update client")
		return
	}
}

func Test_UpdateClient_AuthorizedUser(t *testing.T) {
	reader := bytes.NewReader([]byte("{\"name\": \"updated-name\"}"))

	client, user := models.Common{ID: 1002}, models.Common{ID: 2002}

	context := testutils.New("PATCH", "clients/:id", fmt.Sprintf("/clients/%d", client.ID), "application/json", reader)

	defer context.Database.Close()
	defer clean(context.Database)

	info := struct {
		name  string
		email string
	}{"test1", "test@test.com"}

	context.Request.Client = models.Client{
		Common:       client,
		Name:         "clients-test1",
		ClientID:     "test1-id",
		ClientSecret: "test1-secret",
	}

	context.Request.User = models.User{
		Common: user,
		Name:   &info.name,
		Email:  &info.email,
	}

	context.Database.Create(&context.Request.Client)
	context.Database.Create(&context.Request.User)

	if e := context.Database.Create(&models.ClientAdmin{Client: client.ID, User: user.ID}).Error; e != nil {
		t.Fatalf(fmt.Sprintf("couldn't create client admin fixture: %s", e.Error()))
		return
	}

	err := UpdateClient(&context.Request)

	if err != nil {
		t.Fatalf("authorized user unable to update client: %s", err.Error())
		return
	}
}

func Test_UpdateClient_RandomClientAuthorizedUser(t *testing.T) {
	reader := bytes.NewReader([]byte("{\"name\": \"updated-name\"}"))
	client, user, target := models.Common{ID: 1002}, models.Common{ID: 2002}, models.Common{ID: 3002}
	context := testutils.New("PATCH", "clients/:id", fmt.Sprintf("/clients/%d", target.ID), "application/json", reader)

	defer context.Database.Close()
	defer clean(context.Database)

	info := struct {
		name  string
		email string
	}{"test1", "test@test.com"}

	context.Database.Create(&models.Client{
		Common:       client,
		Name:         "clients-test1",
		ClientID:     "test1-id",
		ClientSecret: "test1-secret",
	})

	context.Database.Create(&models.Client{
		Common:       target,
		Name:         "clients",
		ClientID:     "test2-id",
		ClientSecret: "test2-secret",
	})

	context.Database.Create(&models.User{
		Common: user,
		Name:   &info.name,
		Email:  &info.email,
	})

	context.Request.User = models.User{Common: user}
	context.Request.Client = models.Client{Common: client}

	context.Database.Create(&models.ClientAdmin{Client: target.ID, User: user.ID})

	err := UpdateClient(&context.Request)

	if err != nil {
		t.Fatalf("authorized user unable to update client: %s", err.Error())
		return
	}
}

func Test_UpdateClient_UnauthorizedUser(t *testing.T) {
	reader := bytes.NewReader([]byte("{\"name\": \"updated-name\"}"))
	client, user := models.Common{ID: 1005}, models.Common{ID: 2005}
	context := testutils.New("PATCH", "clients/:id", fmt.Sprintf("/clients/%d", client.ID), "application/json", reader)

	defer context.Database.Close()
	defer clean(context.Database)

	info := struct {
		name  string
		email string
	}{"test1", "test@test.com"}

	context.Request.Client = models.Client{
		Common:       client,
		Name:         "clients-test1",
		ClientID:     "test1-id",
		ClientSecret: "test1-secret",
	}

	context.Request.User = models.User{
		Common: user,
		Name:   &info.name,
		Email:  &info.email,
	}

	context.Database.Create(&context.Request.Client)
	context.Database.Create(&context.Request.User)

	err := UpdateClient(&context.Request)

	if err == nil {
		t.Fatalf("User w/o God privileges was able to update client w/o access")
		return
	}
}
