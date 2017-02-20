package routes

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

func TestUpdateClientRouteGodUser(t *testing.T) {
	body := []byte("{\"name\": \"updated-name\"}")
	reader := bytes.NewReader(body)
	context := testutils.New("PATCH", "clients/:id", "/clients/1337", "application/json", reader)
	defer context.Database.Close()

	db := context.Database

	info := struct {
		name  string
		email string
	}{"test1", "test@test.com"}

	context.Request.Client = models.Client{
		Common:       models.Common{ID: 1337},
		Name:         "clients-test1",
		ClientID:     "123123123",
		ClientSecret: "client-admin-routes-secret",
		System:       true,
	}

	context.Request.User = models.User{
		Common: models.Common{ID: 9999},
		Name:   &info.name,
		Email:  &info.email,
	}

	db.Create(&context.Request.Client)
	db.Create(&context.Request.User)
	mapping := models.UserRoleMapping{User: context.Request.User.ID, Role: 1}
	db.Create(&mapping)

	defer clean(db)

	err := UpdateClient(&context.Request)

	if err != nil {
		t.Fatalf("User w/ God privileges was unable to update client")
		return
	}

	t.Log("successfully updated client w/ God privs")
}

func TestUpdateClientRouteAuthorizedUser(t *testing.T) {
	body := []byte("{\"name\": \"updated-name\"}")
	reader := bytes.NewReader(body)
	context := testutils.New("PATCH", "clients/:id", "/clients/1337", "application/json", reader)
	defer context.Database.Close()

	db := context.Database

	info := struct {
		name  string
		email string
	}{"test1", "test@test.com"}

	context.Request.Client = models.Client{
		Common:       models.Common{ID: 1337},
		Name:         "clients-test1",
		ClientID:     "test1-id",
		ClientSecret: "test1-secret",
	}

	context.Request.User = models.User{
		Common: models.Common{ID: 9999},
		Name:   &info.name,
		Email:  &info.email,
	}

	db.Create(&context.Request.Client)
	db.Create(&context.Request.User)

	mapping := models.ClientAdmin{Client: 1337, User: 9999}
	db.Create(&mapping)

	defer clean(db)

	err := UpdateClient(&context.Request)

	if err != nil {
		t.Fatalf("authorized user unable to update client: %s", err.Error())
		return
	}

	t.Logf("client was updated: %s", context.Request.Client.Name)
}

func TestUpdateClientRouteRandomClientAuthorizedUser(t *testing.T) {
	body := []byte("{\"name\": \"updated-name\"}")
	reader := bytes.NewReader(body)
	context := testutils.New("PATCH", "clients/:id", "/clients/789", "application/json", reader)
	defer context.Database.Close()

	db := context.Database

	info := struct {
		name  string
		email string
	}{"test1", "test@test.com"}

	context.Request.Client = models.Client{
		Common:       models.Common{ID: 456},
		Name:         "clients-test1",
		ClientID:     "test1-id",
		ClientSecret: "test1-secret",
	}

	target := models.Client{
		Common:       models.Common{ID: 789},
		Name:         "clients",
		ClientID:     "test2-id",
		ClientSecret: "test2-secret",
	}

	context.Request.User = models.User{
		Common: models.Common{ID: 123},
		Name:   &info.name,
		Email:  &info.email,
	}

	db.Create(&target)
	db.Create(&context.Request.Client)
	db.Create(&context.Request.User)

	mapping := models.ClientAdmin{Client: 789, User: 123}
	db.Create(&mapping)

	defer clean(db)

	err := UpdateClient(&context.Request)

	if err != nil {
		t.Fatalf("authorized user unable to update client: %s", err.Error())
		return
	}

	t.Logf("client was updated: %s", context.Request.Client.Name)
}

func TestUpdateClientRouteUnauthorizedUser(t *testing.T) {
	body := []byte("{\"name\": \"updated-name\"}")
	reader := bytes.NewReader(body)
	context := testutils.New("PATCH", "clients/:id", "/clients/1337", "application/json", reader)
	defer context.Database.Close()

	db := context.Database

	info := struct {
		name  string
		email string
	}{"test1", "test@test.com"}

	context.Request.Client = models.Client{
		Common:       models.Common{ID: 1337},
		Name:         "clients-test1",
		ClientID:     "test1-id",
		ClientSecret: "test1-secret",
	}

	context.Request.User = models.User{
		Common: models.Common{ID: 9999},
		Name:   &info.name,
		Email:  &info.email,
	}

	db.Create(&context.Request.Client)
	db.Create(&context.Request.User)

	clean := func() {
		db.Exec("DELETE FROM client_admins where id > 1")
		db.Exec("DELETE FROM clients where id > 1")
		db.Exec("DELETE FROM users where id > 1")
		db.Exec("DELETE FROM client_tokens where id > 1")
	}

	defer clean()

	err := UpdateClient(&context.Request)

	if err == nil {
		t.Fatalf("User w/o God privileges was able to update client w/o access")
		return
	}

	t.Logf("successfully blocked attempt to update client w/o access")
}
