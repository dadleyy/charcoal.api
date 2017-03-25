package services

import "testing"
import "github.com/labstack/gommon/log"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/activity"
import "github.com/dadleyy/charcoal.api/testutils"

func Test_Services_Games_IsMember_True(t *testing.T) {
	db := testutils.NewDB()

	name, email := "games-is-member-true", "games-is-member-true@sizethree.cc"
	game, user, member := models.Game{
		Name:   name,
		Status: "ACTIVE",
	}, models.User{
		Email: email,
	}, models.GameMembership{
		Status: "ACTIVE",
	}

	owner := models.User{Email: "game-add-user-owner-2@sizethree.cc"}
	db.Create(&owner)
	defer db.Unscoped().Delete(&owner)

	game.OwnerID = owner.ID

	db.Create(&game)
	defer db.Unscoped().Delete(&game)

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	member.GameID = game.ID
	member.UserID = user.ID
	db.Create(&member)
	defer db.Unscoped().Delete(&member)

	logger := log.New("miritos")
	streams := make(map[string]chan<- activity.Message)
	manager := GameManager{db, logger, streams, game}

	if manager.IsMember(user) != true {
		t.Fatalf("did not recognize valid user/member as a member")
	}
}

func Test_Services_Games_IsMember_False(t *testing.T) {
	db := testutils.NewDB()

	name, email := "games-is-member-false", "games-is-member-false@sizethree.cc"
	game, user := models.Game{Name: name, Status: "ACTIVE"}, models.User{Email: email}

	owner := models.User{Email: "game-add-user-owner-3@sizethree.cc"}
	db.Create(&owner)
	defer db.Unscoped().Delete(&owner)

	game.OwnerID = owner.ID

	db.Create(&game)
	defer db.Unscoped().Delete(&game)

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	logger := log.New("miritos")
	streams := make(map[string]chan<- activity.Message)
	manager := GameManager{db, logger, streams, game}

	if manager.IsMember(user) == true {
		t.Fatalf("did not recognize valid user/member as a member")
	}
}
