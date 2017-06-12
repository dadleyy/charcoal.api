package routes

// import "fmt"
import "bytes"
import "testing"
import "net/url"
import "encoding/json"

import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testutils"
import "github.com/dadleyy/charcoal.api/routes/routetesting"

func createTestRequestBuffer(t interface{}) *bytes.Buffer {
	b, _ := json.Marshal(t)
	return bytes.NewBuffer(b)
}

func Test_Routes_Users_CreateUser_Save(t *testing.T) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Username string `json:"username"`
	}{"dope-1@charcoal.sizethree.cc", "password123", "thename", "user-test-1"}
	clientName := "users-create-1"

	context := routetesting.NewPost(&routetesting.TestRouteParams{}, createTestRequestBuffer(body))

	testutils.CreateClient(&context.Request.Client, clientName, true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	result := CreateUser(context.Request)

	if len(result.Errors) >= 1 {
		t.Fatalf("received error while saving valid user: %v", result)
		return
	}

	client, user, token := context.Request.Client, models.User{}, models.ClientToken{}

	if e := context.Database.Where("email = ?", body.Email).First(&user).Error; e != nil {
		t.Fatalf("unable to find user with email %s", e.Error())
		return
	}

	defer context.Database.Unscoped().Delete(&user)

	if e := context.Database.Where("user_id = ? AND client_id = ?", user.ID, client.ID).First(&token).Error; e != nil {
		t.Fatalf("no token was generated matching user[%d] client[%d] - %s", user.ID, client.ID, e.Error())
		return
	}

	defer context.Database.Unscoped().Delete(&token)
}

func Test_Routes_Users_CreateUser_BadPassword(t *testing.T) {
	body := createTestRequestBuffer(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}{"dope-2@charcoal.sizethree.cc", "password 123", "thename"})

	params := routetesting.TestRouteParams{Values: make(url.Values)}
	context := routetesting.NewPost(&params, body)

	testutils.CreateClient(&context.Request.Client, "bad-password", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	result := CreateUser(context.Request)

	if result == nil || len(result.Errors) == 0 {
		t.Fatalf("should have received error due to bad password")
	}
}

func Test_Routes_Users_CreateUser_DuplicateUsername(t *testing.T) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Username string `json:"username"`
	}{"dope-1@charcoal.sizethree.cc", "password123", "thename", "user-test-1"}
	one := models.User{Email: body.Email + ".diff", Username: body.Username}

	context := routetesting.NewPost(&routetesting.TestRouteParams{}, createTestRequestBuffer(body))

	testutils.CreateClient(&context.Request.Client, "dupe-username", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	context.Database.Create(&one)
	defer context.Database.Unscoped().Delete(&one)

	result := CreateUser(context.Request)

	if result == nil || len(result.Errors) == 0 {
		t.Fatalf("should have received error due to duplicate")
	}
}

/*

func Test_Routes_Users_CreateUser_DuplicateEmail(t *testing.T) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Username string `json:"username"`
	}{"dope-1@charcoal.sizethree.cc", "password123", "thename", "user-test-1"}
	one := models.User{Email: body.Email, Username: body.Username + "-diff"}

	context := routetesting.NewPost(&routetesting.TestRouteParams{}, createTestRequestBuffer(body))

	testutils.CreateClient(&context.Request.Client, "dupe-email", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	context.Database.Create(&one)
	defer context.Database.Unscoped().Delete(&one)

	result := CreateUser(context.Request)

	if result == nil || len(result.Errors) == 0 {
		t.Fatalf("should have received error due to duplicate")
	}
}

func Test_Routes_Users_CreateUser_BadUsername(t *testing.T) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Username string `json:"username"`
	}{"dope-1@charcoal.sizethree.cc", "password123", "thename", "user @ test-1"}

	context := routetesting.NewPost(&routetesting.TestRouteParams{}, createTestRequestBuffer(body))

	testutils.CreateClient(&context.Request.Client, "bad-username", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	result := CreateUser(context.Request)

	if result == nil || len(result.Errors) == 0 {
		t.Fatalf("should have received error due to duplicate")
	}
}

func Test_Routes_Users_UpdateUser_Unauthorized(t *testing.T) {
	db := testutils.NewDB()
	db.DB().SetMaxOpenConns(1)

	user := models.User{
		Email:    "user-update-1@charcoal.sizethree.cc",
		Name:     "thename",
		Username: "user-update-1",
	}

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	body := struct {
		Username string `json:"username"`
	}{"user-update-1-1"}

	params := routetesting.TestRouteParams{}
	params.Set("id", user.ID)

	context := routetesting.NewPatch(&params, createTestRequestBuffer(body))

	testutils.CreateClient(&context.Request.Client, "update-username", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	result := UpdateUser(context.Request)

	if result == nil || len(result.Errors) == 0 {
		t.Fatalf("should have received error due to invalid access")
	}
}

func Test_Routes_Users_UpdateUser_GoodUsername(t *testing.T) {
	db := testutils.NewDB()
	db.DB().SetMaxOpenConns(1)

	user := models.User{
		Email:    "user-update-2@charcoal.sizethree.cc",
		Name:     "thename",
		Username: "user-update-2",
	}

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	body := struct {
		Username string `json:"username"`
	}{"user-update-2-2"}

	params := routetesting.TestRouteParams{}
	params.Set("id", user.ID)
	context := routetesting.NewPatch(&params, createTestRequestBuffer(body))

	context.Request.User = user

	testutils.CreateClient(&context.Request.Client, "update-username", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	result := UpdateUser(context.Request)

	if result != nil && len(result.Errors) >= 1 {
		t.Fatalf("should have been able to update user: %v", result)
	}
}

func Test_Routes_Users_UpdateUser_BadUsername(t *testing.T) {
	db := testutils.NewDB()
	db.DB().SetMaxOpenConns(1)

	user := models.User{
		Email:    "user-update-3@charcoal.sizethree.cc",
		Name:     "thename",
		Username: "user-update-3",
	}

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	body := struct {
		Username string `json:"username"`
	}{"user-update-3 with spaces"}

	params := routetesting.TestRouteParams{}
	params.Set("id", user.ID)
	context := routetesting.NewPatch(&params, createTestRequestBuffer(body))

	context.Request.User = user

	testutils.CreateClient(&context.Request.Client, "update-username", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	result := UpdateUser(context.Request)

	if result == nil || len(result.Errors) == 0 {
		t.Fatalf("should NOT have been able to update user w/ bad username")
	}
}

func Test_Routes_Users_UpdateUser_DuplicateUsername(t *testing.T) {
	db := testutils.NewDB()
	db.DB().SetMaxOpenConns(1)

	one, two := models.User{
		Email:    "user-update-4-1@charcoal.sizethree.cc",
		Name:     "thename",
		Username: "user-update-4-1",
	}, models.User{
		Email:    "user-update-4-2@charcoal.sizethree.cc",
		Name:     "thename",
		Username: "user-update-4-2",
	}

	db.Create(&one)
	defer db.Unscoped().Delete(&one)

	db.Create(&two)
	defer db.Unscoped().Delete(&two)

	body := struct {
		Username string `json:"username"`
	}{"user-update-4-1"}

	params := routetesting.TestRouteParams{}
	params.Set("id", user.ID)
	context := routetesting.NewPatch(&params, createTestRequestBuffer(body))

	context.Request.User = two

	testutils.CreateClient(&context.Request.Client, "update-username-4", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	result := UpdateUser(context.Request)

	if result == nil || len(result.Errors) == 0 {
		t.Fatalf("should NOT have been able to update user w/ duplicate username")
	}
}

func Test_Routes_Users_UpdateUser_DuplicateEmail(t *testing.T) {
	db := testutils.NewDB()
	db.DB().SetMaxOpenConns(1)

	one, two := models.User{
		Email:    "user-update-5-1@charcoal.sizethree.cc",
		Name:     "thename",
		Username: "user-update-5-1",
	}, models.User{
		Email:    "user-update-5-2@charcoal.sizethree.cc",
		Name:     "thename",
		Username: "user-update-5-2",
	}

	db.Create(&one)
	defer db.Unscoped().Delete(&one)

	db.Create(&two)
	defer db.Unscoped().Delete(&two)

	body := struct {
		Email string `json:"email"`
	}{one.Email}

	params := routetesting.TestRouteParams{}
	params.Set("id", user.ID)
	context := routetesting.NewPatch(&params, createTestRequestBuffer(body))

	context.Request.User = two

	testutils.CreateClient(&context.Request.Client, "update-username-5", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	result := UpdateUser(context.Request)

	if result == nil || len(result.Errors) == 0 {
		t.Fatalf("should NOT have been able to update user w/ duplicate username")
	}
}

*/
