package routes

import "testing"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testutils"
import "github.com/dadleyy/charcoal.api/routes/routetesting"

func Test_Routes_ClientAdmins_FindClientAdmins_BadUser_And_NoClient(t *testing.T) {
	ctx := routetesting.NewFind("client-admins")

	if err := FindClientAdmins(&ctx.Request); err != nil {
		return
	}

	t.Fatalf("should not have passed w/o error")
}

func Test_Routes_ClientAdmins_FindClientAdmins_BadUser_With_Client(t *testing.T) {
	db := testutils.NewDB()

	client := models.Client{}
	testutils.CreateClient(&client, "client-admins-find-1", false)
	defer db.Unscoped().Delete(&client)

	ctx := routetesting.NewFind("client-admins")
	ctx.Request.Client = client

	if err := FindClientAdmins(&ctx.Request); err != nil {
		return
	}

	t.Fatalf("should not have passed w/o error")
}

func Test_Routes_ClientAdmins_FindClientAdmins_ValidClientAdmin(t *testing.T) {
	db := testutils.NewDB()

	email := "client-admins-find-2@charcoal.sizethree.cc"
	client, user := models.Client{}, models.User{Email: email}

	testutils.CreateClient(&client, "client-admins-find-2", false)
	defer db.Unscoped().Delete(&client)

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	mapping := models.ClientAdmin{Client: client.ID, User: user.ID}
	db.Create(&mapping)
	defer db.Unscoped().Delete(&mapping)

	ctx := routetesting.NewFind("client-admins")
	ctx.Request.Client = client
	ctx.Request.User = user

	if err := FindClientAdmins(&ctx.Request); err != nil {
		t.Fatalf("error even though user is admin: %s", err.Error())
		return
	}
}

func Test_Routes_ClientAdmins_FindClientAdmins_ValidGodUser(t *testing.T) {
	db := testutils.NewDB()

	email := "client-admins-find-3@charcoal.sizethree.cc"
	client, user := models.Client{}, models.User{Email: email}

	testutils.CreateClient(&client, "client-admins-find-2", false)
	defer db.Unscoped().Delete(&client)

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	mapping := models.UserRoleMapping{User: user.ID, Role: 1}
	db.Create(&mapping)
	defer db.Unscoped().Delete(&mapping)

	ctx := routetesting.NewFind("client-admins")
	ctx.Request.Client = client
	ctx.Request.User = user

	if err := FindClientAdmins(&ctx.Request); err != nil {
		t.Fatalf("error even though user is admin: %s", err.Error())
		return
	}
}
