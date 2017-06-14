package routes

import "testing"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testing/utils"
import "github.com/dadleyy/charcoal.api/testing/routing"

func Test_Routes_ClientAdmins_FindClientAdmins_BadUser_And_NoClient(t *testing.T) {
	ctx := testrouting.NewFind(&testrouting.TestRouteParams{})

	r := FindClientAdmins(ctx.Request)

	if len(r.Errors) == 0 {
		t.Fatalf("should not have passed w/o error")
	}
}

func Test_Routes_ClientAdmins_FindClientAdmins_BadUser_With_Client(t *testing.T) {
	db := testutils.NewDB()

	client := models.Client{}
	testutils.CreateClient(&client, "client-admins-find-1", false)
	defer db.Unscoped().Delete(&client)

	ctx := testrouting.NewFind(&testrouting.TestRouteParams{})
	ctx.Request.Client = client

	r := FindClientAdmins(ctx.Request)

	if len(r.Errors) == 0 {
		t.Fatalf("should not have passed w/o error")
	}
}

func Test_Routes_ClientAdmins_FindClientAdmins_ValidClientAdmin(t *testing.T) {
	db := testutils.NewDB()

	email := "client-admins-find-2@charcoal.sizethree.cc"
	client, user := models.Client{}, models.User{Email: email}

	testutils.CreateClient(&client, "client-admins-find-2", false)
	defer db.Unscoped().Delete(&client)

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	mapping := models.ClientAdmin{ClientID: client.ID, UserID: user.ID}
	db.Create(&mapping)
	defer db.Unscoped().Delete(&mapping)

	ctx := testrouting.NewFind(&testrouting.TestRouteParams{})
	ctx.Request.Client = client
	ctx.Request.User = user

	r := FindClientAdmins(ctx.Request)

	if len(r.Errors) >= 1 {
		t.Fatalf("error even though user is admin")
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

	mapping := models.UserRoleMapping{UserID: user.ID, RoleID: 1}
	db.Create(&mapping)
	defer db.Unscoped().Delete(&mapping)

	ctx := testrouting.NewFind(&testrouting.TestRouteParams{})
	ctx.Request.Client = client
	ctx.Request.User = user

	r := FindClientAdmins(ctx.Request)

	if len(r.Errors) >= 1 {
		t.Fatalf("error even though user is admin: %s", r.Errors[0].Error())
	}
}
