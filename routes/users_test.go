package routes

import "fmt"
import "bytes"
import "testing"
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

	context := routetesting.NewPost("users", createTestRequestBuffer(body))

	testutils.CreateClient(&context.Request.Client, clientName, true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	err := CreateUser(&context.Request)

	if err != nil {
		t.Fatalf("received error while saving valid user: %s", err.Error())
		return
	}

	client, user, token := context.Request.Client, models.User{}, models.ClientToken{}

	if e := context.Database.Where("email = ?", body.Email).First(&user).Error; e != nil {
		t.Fatalf("unable to find user with email %s", e.Error())
		return
	}

	defer context.Database.Unscoped().Delete(&user)

	if e := context.Database.Where("user = ? AND client = ?", user.ID, client.ID).First(&token).Error; e != nil {
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

	context := routetesting.NewPost("users", body)

	testutils.CreateClient(&context.Request.Client, "bad-password", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	testutils.CreateClient(&context.Request.Client, "users_create_client", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	err := CreateUser(&context.Request)

	if err == nil {
		t.Fatalf("should have received error due to bad password")
		return
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

	context := routetesting.NewPost("users", createTestRequestBuffer(body))

	testutils.CreateClient(&context.Request.Client, "dupe-username", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	context.Database.Create(&one)
	defer context.Database.Unscoped().Delete(&one)

	err := CreateUser(&context.Request)

	if err == nil {
		t.Fatalf("should have received error due to duplicate")
		return
	}
}

func Test_Routes_Users_CreateUser_DuplicateEmail(t *testing.T) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Username string `json:"username"`
	}{"dope-1@charcoal.sizethree.cc", "password123", "thename", "user-test-1"}
	one := models.User{Email: body.Email, Username: body.Username + "-diff"}

	context := routetesting.NewPost("users", createTestRequestBuffer(body))

	testutils.CreateClient(&context.Request.Client, "dupe-email", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	context.Database.Create(&one)
	defer context.Database.Unscoped().Delete(&one)

	err := CreateUser(&context.Request)

	if err == nil {
		t.Fatalf("should have received error due to duplicate")
		return
	}
}

func Test_Routes_Users_CreateUser_BadUsername(t *testing.T) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Username string `json:"username"`
	}{"dope-1@charcoal.sizethree.cc", "password123", "thename", "user @ test-1"}

	context := routetesting.NewPost("users", createTestRequestBuffer(body))

	testutils.CreateClient(&context.Request.Client, "bad-username", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	err := CreateUser(&context.Request)

	if err == nil {
		t.Fatalf("should have received error due to duplicate")
		return
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

	context := routetesting.NewPatch("users/:id", fmt.Sprintf("users/%d", user.ID), createTestRequestBuffer(body))

	testutils.CreateClient(&context.Request.Client, "update-username", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	err := UpdateUser(&context.Request)

	if err == nil {
		t.Fatalf("should NOT have been able to update user - no user associated w/ request")
		return
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

	context := routetesting.NewPatch("users/:id", fmt.Sprintf("users/%d", user.ID), createTestRequestBuffer(body))

	context.Request.User = user

	testutils.CreateClient(&context.Request.Client, "update-username", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	err := UpdateUser(&context.Request)

	if err != nil {
		t.Fatalf("should have been able to update user: %s", err.Error())
		return
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

	context := routetesting.NewPatch("users/:id", fmt.Sprintf("users/%d", user.ID), createTestRequestBuffer(body))

	context.Request.User = user

	testutils.CreateClient(&context.Request.Client, "update-username", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	err := UpdateUser(&context.Request)

	if err == nil {
		t.Fatalf("should NOT have been able to update user w/ bad username")
		return
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

	context := routetesting.NewPatch("users/:id", fmt.Sprintf("users/%d", two.ID), createTestRequestBuffer(body))

	context.Request.User = two

	testutils.CreateClient(&context.Request.Client, "update-username-4", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	err := UpdateUser(&context.Request)

	if err != nil {
		return
	}

	t.Fatalf("should NOT have been able to update user w/ duplicate username")
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

	context := routetesting.NewPatch("users/:id", fmt.Sprintf("users/%d", two.ID), createTestRequestBuffer(body))

	context.Request.User = two

	testutils.CreateClient(&context.Request.Client, "update-username-5", true)
	defer context.Database.Unscoped().Delete(&context.Request.Client)

	err := UpdateUser(&context.Request)

	if err != nil {
		return
	}

	t.Fatalf("should NOT have been able to update user w/ duplicate email")
}
