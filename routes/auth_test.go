package routes

import "fmt"
import "bytes"
import "testing"
import "golang.org/x/crypto/bcrypt"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testutils"
import "github.com/dadleyy/charcoal.api/routes/routetesting"

func Test_Routes_Auth_PasswordLogin_NoClient(t *testing.T) {
	body := bytes.NewBuffer([]byte("{}"))
	ctx := routetesting.NewPost(&routetesting.TestRouteParams{}, body)

	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	user := models.User{Email: "login-test-1@sizethree.cc", Password: string(password)}
	client := models.Client{}
	token := models.ClientToken{}

	ctx.Database.Create(&user)
	defer ctx.Database.Unscoped().Delete(&user)

	testutils.CreateClient(&client, "login-test-1", true)
	defer ctx.Database.Unscoped().Delete(&client)

	token.ClientID = client.ID
	token.UserID = user.ID
	ctx.Database.Create(&token)
	defer ctx.Database.Unscoped().Delete(&token)

	if err := PasswordLogin(ctx.Request); err != nil {
		return
	}

	t.Fatalf("No error but client was not set on request runtime!")
}

func Test_Routes_Auth_PasswordLogin_NonSystem(t *testing.T) {
	body := bytes.NewBuffer([]byte("{}"))
	ctx := routetesting.NewPost(&routetesting.TestRouteParams{}, body)

	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	user := models.User{Email: "login-test-2@sizethree.cc", Password: string(password)}
	client := models.Client{}
	token := models.ClientToken{}

	ctx.Database.Create(&user)
	defer ctx.Database.Unscoped().Delete(&user)

	testutils.CreateClient(&client, "login-test-2", false)
	defer ctx.Database.Unscoped().Delete(&client)

	token.ClientID = client.ID
	token.UserID = user.ID
	ctx.Database.Create(&token)
	defer ctx.Database.Unscoped().Delete(&token)

	ctx.Request.Client = client
	r := PasswordLogin(ctx.Request)

	if len(r.Errors) == 0 {
		t.Fatalf("No error but client was not set on request runtime!")
	}
}

func Test_Routes_Auth_PasswordLogin_SystemGoodPassword(t *testing.T) {
	email, password := "lt2@sizethree.cc", "password"
	body := bytes.NewBuffer([]byte(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", email, password)))

	ctx := routetesting.NewPost(&routetesting.TestRouteParams{}, body)

	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Email: email, Password: string(pw)}

	client := models.Client{}
	token := models.ClientToken{}

	ctx.Database.Create(&user)
	defer ctx.Database.Unscoped().Delete(&user)

	testutils.CreateClient(&client, "login-test-2", true)
	defer ctx.Database.Unscoped().Delete(&client)

	token.ClientID = client.ID
	token.UserID = user.ID
	ctx.Database.Create(&token)
	defer ctx.Database.Unscoped().Delete(&token)

	ctx.Request.Client = client
	r := PasswordLogin(ctx.Request)

	if len(r.Errors) >= 1 {
		t.Fatalf("valid login failed: %v", r)
	}
}

func Test_Routes_Auth_PasswordLogin_SystemBadPassword(t *testing.T) {
	email, password := "lt2@sizethree.cc", "password"
	body := bytes.NewBuffer([]byte(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"fudge\"}", email)))
	ctx := routetesting.NewPost(&routetesting.TestRouteParams{}, body)

	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Email: email, Password: string(pw)}

	client := models.Client{}
	token := models.ClientToken{}

	ctx.Database.Create(&user)
	defer ctx.Database.Unscoped().Delete(&user)

	testutils.CreateClient(&client, "login-test-2", true)
	defer ctx.Database.Unscoped().Delete(&client)

	token.ClientID = client.ID
	token.UserID = user.ID
	ctx.Database.Create(&token)
	defer ctx.Database.Unscoped().Delete(&token)

	ctx.Request.Client = client
	r := PasswordLogin(ctx.Request)

	if len(r.Errors) == 0 {
		t.Fatalf("invalid password should have failed.")
	}
}
