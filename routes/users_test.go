package routes

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
	body := createTestRequestBuffer(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}{"dope@charcoal.sizethree.cc", "password123", "thename"})
	return

	context := routetesting.NewPost("users", body)
	defer context.Database.Close()

	testutils.CreateClient(&context.Request.Client, "users_create_client", true)

	err := CreateUser(&context.Request)

	if err != nil {
		t.Fatalf("received error while saving valid user: %s", err.Error())
		return
	}

	client, user, token := models.Client{}, models.User{}, models.ClientToken{}

	if e := context.Database.Where("email = ?", "dope@charcoal.sizethree.cc").First(&user).Error; e != nil {
		t.Fatalf("unable to find user with email %s", e.Error())
		return
	}

	if e := context.Database.Where("name = ?", "users_create_client").First(&client).Error; e != nil {
		t.Fatalf("unable to find client matching name: %s", e.Error())
		return
	}

	if e := context.Database.Where("user = ? AND client = ?", user.ID, client.ID).First(&token).Error; e == nil {
		return
	}

	t.Fatalf("no token was generated matching user[%d] client[%d]", user.ID, client.ID)
}
