package services

import "testing"
import "github.com/dadleyy/charcoal.api/util"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testutils"

func Test_Services_UserClients_Associate(t *testing.T) {
	db := testutils.NewDB()

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

	defer db.Unscoped().Delete(&client)
	defer db.Unscoped().Delete(&user)

	var c int

	db.Model(&models.ClientToken{}).Where("user_id = ? AND client_id = ?", user.ID, client.ID).Count(&c)

	if c != 0 {
		t.Fatalf("newly created client and user already associated?")
		return
	}

	if _, err := mgr.Associate(&user, &client); err != nil {
		t.Fatalf("failed associate: %s", err.Error())
		return
	}

	cursor := db.Model(&models.ClientToken{}).Where("user_id = ? AND client_id = ?", user.ID, client.ID)
	defer db.Unscoped().Where("user_id = ? AND client_id = ?", user.ID, client.ID).Delete(models.ClientToken{})

	if _ = cursor.Count(&c).Error; c == 1 {
		return
	}
	defer db.Unscoped().Delete(models.ClientToken{})

	t.Fatalf("unable to find newly associted token")
}
