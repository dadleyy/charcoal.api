package services

import "testing"
import "github.com/labstack/gommon/log"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/activity"
import "github.com/dadleyy/charcoal.api/testutils"

func Test_Services_Games_AddUser_EntryID(t *testing.T) {
	db := testutils.NewDB()

	name, email := "games-add-user-entry-id-1", "games-add-user-entry-id-1@sizethree.cc"
	game, round, user := models.Game{Name: name, Status: "ACTIVE"}, models.GameRound{}, models.User{Email: email}

	owner := models.User{Email: "game-add-user-owner-1@sizethree.cc"}
	db.Create(&owner)
	defer db.Unscoped().Delete(&owner)

	game.OwnerID = owner.ID

	db.Create(&game)
	defer db.Unscoped().Delete(&game)

	round.GameID = game.ID
	db.Create(&round)
	defer db.Unscoped().Delete(&round)

	db.Create(&user)
	defer db.Unscoped().Delete(&user)

	logger := log.New("miritos")

	streams := make(map[string]chan<- activity.Message)
	manager := GameManager{db, logger, streams, game}
	manager.AddUser(user)

	member := models.GameMembership{}

	if e := db.Where("game_id = ?", game.ID).First(&member).Error; e != nil {
		t.Fatalf("user was not added to game!")
		return
	}

	defer db.Unscoped().Delete(&member)

	if member.EntryRoundID == nil || *member.EntryRoundID != round.ID {
		t.Fatalf("entry round id was nil when it should have been %d", round.ID)
	}
}

func Test_Services_Games_AddUser_NoEntryID(t *testing.T) {
	db := testutils.NewDB()

	name, email := "games-add-user-entry-id-2", "games-add-user-entry-id-2@sizethree.cc"
	game, user := models.Game{Name: name, Status: "ACTIVE"}, models.User{Email: email}

	owner := models.User{Email: "game-add-user-owner-4@sizethree.cc"}
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
	manager.AddUser(user)

	member := models.GameMembership{}

	if e := db.Where("game_id = ?", game.ID).First(&member).Error; e != nil {
		t.Fatalf("user was not added to game!")
		return
	}

	defer db.Unscoped().Delete(&member)

	if member.EntryRoundID != nil {
		t.Fatalf("entry round id was not when it should have been")
	}
}

func Test_Services_Games_IsMember_True(t *testing.T) {
	db := testutils.NewDB()

	name, email := "games-is-member-true", "games-is-member-true@sizethree.cc"
	game, user, member := models.Game{Name: name, Status: "ACTIVE"}, models.User{Email: email}, models.GameMembership{}

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
