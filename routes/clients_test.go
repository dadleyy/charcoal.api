package routes

import "fmt"
import "bytes"
import "testing"

import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testutils"
import "github.com/dadleyy/charcoal.api/routes/routetesting"

func Test_Routes_Clients_UpdateClient_GodUser(t *testing.T) {
	reader := bytes.NewReader([]byte("{\"name\": \"updated-name\"}"))

	email := "clients-test-1@charcoal.sizethree.cc"

	db := testutils.NewDB()

	client, user := models.Client{}, models.User{Email: email}

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	mapping := models.UserRoleMapping{User: user.ID, Role: 1}
	db.Create(&mapping)
	defer db.Unscoped().Delete(&mapping)

	testutils.CreateClient(&client, "clients-test-1", true)
	defer db.Unscoped().Delete(&client)

	ctx := routetesting.NewPatch("clients/:id", fmt.Sprintf("/clients/%d", client.ID), reader)

	ctx.Request.Client = client
	ctx.Request.User = user

	result := UpdateClient(ctx.Request)

	if len(result.Errors) >= 1 {
		t.Fatalf("god user was unable to update client: %v", result)
	}
}

func Test_Routes_Clients_UpdateClient_AuthorizedUser(t *testing.T) {
	reader := bytes.NewReader([]byte("{\"name\": \"updated-name\"}"))

	email := "clients-test-2@charcoal.sizethree.cc"
	db := testutils.NewDB()
	client, user := models.Client{}, models.User{Email: email}

	db.Create(&user)

	defer db.Unscoped().Delete(&user)

	testutils.CreateClient(&client, "clients-test-2", true)
	defer db.Unscoped().Delete(&client)

	mapping := models.ClientAdmin{UserID: user.ID, ClientID: client.ID}
	db.Create(&mapping)
	defer db.Unscoped().Delete(&mapping)

	ctx := routetesting.NewPatch("clients/:id", fmt.Sprintf("/clients/%d", client.ID), reader)

	ctx.Request.Client = client
	ctx.Request.User = user

	result := UpdateClient(ctx.Request)

	if len(result.Errors) >= 1 {
		t.Fatalf("god user was unable to update client: %v", result)
	}
}

func Test_Routes_Clients_UpdateClient_OtherClientAuthorizedUser(t *testing.T) {
	reader := bytes.NewReader([]byte("{\"name\": \"updated-name\"}"))
	email := "clients-test-3@charcoal.sizethree.cc"
	db := testutils.NewDB()
	client, target, user := models.Client{}, models.Client{}, models.User{Email: email}

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	testutils.CreateClient(&target, "clients-test-3.1", true)
	defer db.Unscoped().Delete(&target)

	testutils.CreateClient(&client, "clients-test-3.2", true)
	defer db.Unscoped().Delete(&client)

	// associate the user w/ our target client
	mapping := models.ClientAdmin{UserID: user.ID, ClientID: target.ID}
	db.Create(&mapping)
	defer db.Unscoped().Delete(&mapping)

	// use the target client's id in the route param
	ctx := routetesting.NewPatch("clients/:id", fmt.Sprintf("/clients/%d", target.ID), reader)

	// use the other client as the request runtime
	ctx.Request.Client = client
	ctx.Request.User = user

	result := UpdateClient(ctx.Request)

	if len(result.Errors) >= 1 {
		t.Fatalf("client admin was unable to update their client from another: %v", result)
	}
}

func Test_Routes_Clients_UpdateClient_UnauthorizedUser(t *testing.T) {
	reader := bytes.NewReader([]byte("{\"name\": \"updated-name\"}"))

	email := "clients-test-4@charcoal.sizethree.cc"
	db := testutils.NewDB()
	client, user := models.Client{}, models.User{Email: email}

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	testutils.CreateClient(&client, "clients-test-4", true)
	defer db.Unscoped().Delete(&client)

	ctx := routetesting.NewPatch("clients/:id", fmt.Sprintf("/clients/%d", client.ID), reader)

	ctx.Request.Client = client
	ctx.Request.User = user

	result := UpdateClient(ctx.Request)

	if result != nil && len(result.Errors) == 0 {
		t.Fatalf("invalid user able to update client")
	}
}
