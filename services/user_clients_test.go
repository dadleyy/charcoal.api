package services

import "testing"
import "github.com/dadleyy/charcoal.api/util"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testutils"

func Test_Services_UserClients_Associate(t *testing.T) {
	db := testutils.NewDB()

	defer db.Exec("DELETE FROM users where id > 0")
	defer db.Exec("DELETE FROM clients where id > 0")

	defer db.Exec("DELETE FROM user_role_mappings where id > 0")
	defer db.Exec("DELETE FROM client_tokens where id > 0")

	mgr := UserClientManager{db}

	client := models.Client{
		Name:         "whoa buddy",
		ClientID:     util.RandStringBytesMaskImprSrc(20),
		ClientSecret: util.RandStringBytesMaskImprSrc(40),
	}

	email := "test-associate@charcoal.sizethree.cc"
	user := models.User{Email: email}

	db.Create(&client)
	db.Create(&user)

	var c int

	db.Model(&models.ClientToken{}).Where("user = ? AND client = ?", user.ID, client.ID).Count(&c)

	if c != 0 {
		t.Fatalf("newly created client and user already associated?")
		return
	}

	if _, err := mgr.Associate(&user, &client); err != nil {
		t.Fatalf("failed associate: %s", err.Error())
		return
	}

	if _ = db.Model(&models.ClientToken{}).Where("user = ? AND client = ?", user.ID, client.ID).Count(&c).Error; c == 1 {
		return
	}

	t.Fatalf("unable to find newly associted token")
}
