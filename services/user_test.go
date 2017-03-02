package services

import "testing"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testutils"

func Test_Services_Users_IsDuplicateTrue(t *testing.T) {
	db := testutils.NewDB()
	defer db.Close()

	defer db.Exec("DELETE FROM users where id > 0")
	defer db.Exec("DELETE FROM user_role_mappings where id > 0")

	dupe := "testing@charcoal.sizethree.cc"
	db.Create(&models.User{Email: &dupe})
	mgr := UserManager{db}

	if t, _ := mgr.IsDuplicate(&models.User{Email: &dupe}); t == true {
		return
	}

	t.Fatalf("duplicate check failed")
}

func Test_Services_Users_IsDuplicateFalse(t *testing.T) {
	db := testutils.NewDB()
	defer db.Close()

	defer db.Exec("DELETE FROM users where id > 0")
	defer db.Exec("DELETE FROM user_role_mappings where id > 0")

	dupe := "testing@charcoal.sizethree.cc"
	nodupe := "testing-2@charcoal.sizethree.cc"
	db.Create(&models.User{Email: &dupe})
	mgr := UserManager{db}

	if t, _ := mgr.IsDuplicate(&models.User{Email: &nodupe}); t == false {
		return
	}

	t.Fatalf("duplicate check failed")
}

func Test_Services_Users_IsAdminTrue(t *testing.T) {
	db := testutils.NewDB()
	defer db.Close()

	defer db.Exec("DELETE FROM users where id > 0")
	defer db.Exec("DELETE FROM user_role_mappings where id > 0")

	dupe := "testing@charcoal.sizethree.cc"
	user := models.User{Email: &dupe}
	db.Create(&user)
	db.Create(&models.UserRoleMapping{Role: 1, User: user.ID})

	mgr := UserManager{db}

	if t := mgr.IsAdmin(&user); t == true {
		return
	}

	t.Fatalf("duplicate check failed")
}

func Test_Services_Users_IsAdminFalse(t *testing.T) {
	db := testutils.NewDB()
	defer db.Close()

	defer db.Exec("DELETE FROM users where id > 0")
	defer db.Exec("DELETE FROM user_role_mappings where id > 0")

	dupe := "testing@charcoal.sizethree.cc"
	user := models.User{Email: &dupe}
	db.Create(&user)

	mgr := UserManager{db}

	if t := mgr.IsAdmin(&user); t == false {
		return
	}

	t.Fatalf("duplicate check failed")
}
