package services

import "testing"
import "github.com/jinzhu/gorm"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/defs"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/activity"
import "github.com/dadleyy/charcoal.api/testing/utils"

type gameServiceTestContext struct {
	*gorm.DB

	users   []*models.User
	game    *models.Game
	manager *GameManager
}

type gameServiceTestLifecycle struct {
	setup    chan gameServiceTestContext
	teardown chan *testing.T
}

func (life *gameServiceTestLifecycle) done(t *testing.T) {
	life.teardown <- t
}

func (life *gameServiceTestLifecycle) begin() {
	db := testutils.NewDB()

	users := []*models.User{
		&models.User{Email: "add-user-test-new-1@test.com"},
		&models.User{Email: "add-user-test-new-2@test.com"},
	}

	for i, _ := range users {
		db.Create(users[i])
		defer db.Unscoped().Delete(users[i])
	}

	game := models.Game{
		Name:    "test-game",
		Status:  "ACTIVE",
		OwnerID: users[0].ID,
	}

	db.Create(&game)
	defer db.Unscoped().Delete(&game)

	manager := GameManager{
		DB:      db,
		Logger:  log.New("miritos"),
		Streams: make(map[string]chan<- activity.Message),
		Game:    game,
	}

	ctx := gameServiceTestContext{
		DB:      db,
		users:   users,
		game:    &game,
		manager: &manager,
	}

	life.setup <- ctx
	<-life.teardown
}

func newGameServiceTestLife() gameServiceTestLifecycle {
	setup, teardown := make(chan gameServiceTestContext), make(chan *testing.T)
	return gameServiceTestLifecycle{setup, teardown}
}

func Test_Services_Games_AddUser_New(t *testing.T) {
	lf := newGameServiceTestLife()
	go lf.begin()
	ctx := <-lf.setup
	defer lf.done(t)

	membership, e := ctx.manager.AddUser(*ctx.users[1])

	if e != nil {
		t.Fatalf("errored while adding user to game: %s", e.Error())
		return
	}

	defer ctx.Unscoped().Delete(&membership)

	if membership.Status != defs.GameMembershipActiveStatus {
		t.Fatalf("expected new member to have active status but found: %s", membership.Status)
		return
	}
}

func Test_Services_Games_AddUser_Inactive(t *testing.T) {
	lf := newGameServiceTestLife()
	go lf.begin()
	ctx := <-lf.setup
	defer lf.done(t)

	membership := models.GameMembership{
		UserID: ctx.users[1].ID,
		GameID: ctx.game.ID,
		Status: defs.GameMembershipInactiveStatus,
	}

	ctx.Create(&membership)
	defer ctx.Unscoped().Delete(&membership)

	result, e := ctx.manager.AddUser(*ctx.users[1])

	if e != nil {
		t.Fatalf("expected successfull membership addition but receieved: %s", e.Error())
		return
	}

	if result.Status != defs.GameMembershipActiveStatus {
		t.Fatalf("expected active membership status but found: %s", result.Status)
		return
	}

	if result.ID != membership.ID {
		t.Fatalf("expected to re-use membership but ids didnt match: %d != %d", result.ID, membership.ID)
		return
	}
}

func Test_Services_Games_AddUser_Active(t *testing.T) {
	lf := newGameServiceTestLife()
	go lf.begin()
	ctx := <-lf.setup
	defer lf.done(t)

	membership := models.GameMembership{
		UserID: ctx.users[1].ID,
		GameID: ctx.game.ID,
		Status: defs.GameMembershipActiveStatus,
	}

	ctx.Create(&membership)
	defer ctx.Unscoped().Delete(&membership)

	beforeCount, afterCount := -1, -1

	ctx.Model(&models.GameMembership{}).Where("game_id = ?", ctx.game.ID).Count(&beforeCount)

	_, e := ctx.manager.AddUser(*ctx.users[1])

	ctx.Model(&models.GameMembership{}).Where("game_id = ?", ctx.game.ID).Count(&afterCount)

	if e == nil || beforeCount != afterCount {
		t.Fatalf("expected error adding active membership but did not receive one")
		return
	}
}

func Test_Services_Games_RemoveMember_Active(t *testing.T) {
	lf := newGameServiceTestLife()
	go lf.begin()
	ctx := <-lf.setup
	defer lf.done(t)

	membership := models.GameMembership{
		UserID: ctx.users[1].ID,
		GameID: ctx.game.ID,
		Status: defs.GameMembershipActiveStatus,
	}

	ctx.Create(&membership)
	defer ctx.Unscoped().Delete(&membership)

	e := ctx.manager.RemoveMember(membership)

	if e != nil {
		t.Fatalf("unexpected error: %s", e.Error())
		return
	}

	after := models.GameMembership{}

	ctx.First(&after, membership.ID)

	if after.Status != defs.GameMembershipInactiveStatus {
		t.Fatalf("expected inactive status but found: %s", after.Status)
		return
	}
}

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
