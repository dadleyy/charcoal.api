package services

import "testing"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testing/utils"

func Test_Services_Users_IsDuplicateTrue(t *testing.T) {
	db := testutils.NewDB()

	dupe := "testing@charcoal.sizethree.cc"
	user := models.User{Email: dupe}
	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	mgr := UserManager{db, log.New("t")}

	if t, _ := mgr.IsDuplicate(&models.User{Email: dupe}); t == true {
		return
	}

	t.Fatalf("duplicate check failed")
}

func Test_Services_Users_IsDuplicateFalse(t *testing.T) {
	db := testutils.NewDB()

	dupe := "testing@charcoal.sizethree.cc"
	nodupe := "testing-2@charcoal.sizethree.cc"
	user := models.User{Email: dupe}
	db.Create(&user)

	defer db.Unscoped().Delete(&user)

	mgr := UserManager{db, log.New("t")}

	if t, _ := mgr.IsDuplicate(&models.User{Email: nodupe}); t == false {
		return
	}

	t.Fatalf("duplicate check failed")
}

func Test_Services_Users_IsAdminTrue(t *testing.T) {
	db := testutils.NewDB()

	dupe := "testing@charcoal.sizethree.cc"

	user := models.User{Email: dupe}
	db.Create(&user)

	mapping := models.UserRoleMapping{RoleID: 1, UserID: user.ID}
	db.Create(&mapping)

	defer db.Unscoped().Delete(&user)
	defer db.Unscoped().Delete(&mapping)

	mgr := UserManager{db, log.New("t")}

	if t := mgr.IsAdmin(&user); t == true {
		return
	}

	t.Fatalf("duplicate check failed")
}

func Test_Services_Users_IsAdminFalse(t *testing.T) {
	db := testutils.NewDB()

	dupe := "testing@charcoal.sizethree.cc"
	user := models.User{Email: dupe}
	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	mgr := UserManager{db, log.New("t")}

	if t := mgr.IsAdmin(&user); t == false {
		return
	}

	t.Fatalf("duplicate check failed")
}
