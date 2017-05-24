package activity

import "fmt"
import "sync"
import "testing"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/defs"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testutils"

func Test_Activity_GameProcessor_JoinedMessage(t *testing.T) {
	db, logger, stream := testutils.NewDB(), log.New("miritos"), make(chan Message, 1)

	processor := GameProcessor{
		Logger: logger,
		DB:     db,
		Stream: stream,
	}

	email := "game-processor-test-1.1@charcoal.sizethree.cc"
	ownerEmail := "game-processor-test-1.2@charcoal.sizethree.cc"
	game, user, owner := models.Game{Status: "ACTIVE"}, models.User{Email: email}, models.User{Email: ownerEmail}

	db.Create(&owner)
	db.Create(&user)

	game.OwnerID = owner.ID
	db.Create(&game)

	defer db.Unscoped().Delete(&owner)
	defer db.Unscoped().Delete(&user)
	defer db.Unscoped().Delete(&game)

	mem := models.GameMembership{
		UserID: user.ID,
		GameID: game.ID,
		Status: defs.GameMembershipActiveStatus,
	}

	db.Create(&mem)
	defer db.Unscoped().Delete(&mem)
	defer db.Unscoped().Where("game_id = ?", game.ID).Delete(models.GameMembershipHistory{})

	stream <- Message{&user, &game, fmt.Sprintf("%s:%s", defs.GamesStreamIdentifier, defs.GameProcessorUserJoined)}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go processor.Begin(&wg)
	close(stream)
	wg.Wait()

	g := models.Game{}
	db.First(&g, game.ID)

	if g.Population == 1 {
		return
	}

	t.Fatalf("expected 1, found: %d", g.Population)
}

func Test_Activity_GameProcessor_LeftMessage(t *testing.T) {
	db, logger, stream := testutils.NewDB(), log.New("miritos"), make(chan Message, 1)

	processor := GameProcessor{
		Logger: logger,
		DB:     db,
		Stream: stream,
	}

	email := "game-processor-test-2.1@charcoal.sizethree.cc"
	ownerEmail := "game-processor-test-2.2@charcoal.sizethree.cc"
	game, user, owner := models.Game{
		Status:     "ACTIVE",
		Population: 123123,
	}, models.User{
		Email: email,
	}, models.User{
		Email: ownerEmail,
	}

	db.Create(&owner)
	db.Create(&user)

	game.OwnerID = owner.ID
	db.Create(&game)

	defer db.Unscoped().Delete(&owner)
	defer db.Unscoped().Delete(&user)
	defer db.Unscoped().Delete(&game)

	round := models.GameRound{GameID: game.ID}
	db.Create(&round)
	defer db.Unscoped().Delete(&round)

	history := models.GameMembershipHistory{UserID: user.ID, GameID: game.ID}
	db.Create(&history)
	defer db.Unscoped().Where("game_id = ?", game.ID).Delete(history)

	mem := models.GameMembership{
		UserID: user.ID,
		GameID: game.ID,
		Status: defs.GameMembershipActiveStatus,
	}

	db.Create(&mem)
	defer db.Unscoped().Delete(&mem)

	stream <- Message{&user, &game, fmt.Sprintf("%s:%s", defs.GamesStreamIdentifier, defs.GameProcessorUserLeft)}
	wg := sync.WaitGroup{}

	wg.Add(1)
	go processor.Begin(&wg)
	close(stream)
	wg.Wait()

	g := models.Game{}
	db.First(&g, game.ID)

	if g.Population == 1 {
		return
	}

	t.Fatalf("expected 1, found: %d", g.Population)
}
