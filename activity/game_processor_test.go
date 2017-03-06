package activity

import "fmt"
import "testing"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testutils"

func Test_Activity_GameProcessor_JoinedMessage(t *testing.T) {
	db := testutils.NewDB()
	defer db.Close()
	logger := log.New("miritos")

	stream := make(chan Message, 1)
	wait := make(chan struct{})

	processor := GameProcessor{logger, stream, wait, db}

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

	stream <- Message{&user, &game, fmt.Sprintf("games:%s", GameProcessorUserJoined)}
	go processor.Begin(ProcessorConfig{DB: testutils.DBConfig()})
	close(stream)
	<-wait

	g := models.Game{}
	db.First(&g, game.ID)

	if g.Population == 1 {
		return
	}

	t.Fatalf("expected game w/ initial population of 0 to increase to 1 after update")
}

func Test_Activity_GameProcessor_LeftMessage(t *testing.T) {
	db := testutils.NewDB()
	defer db.Close()
	logger := log.New("miritos")

	stream := make(chan Message, 1)
	wait := make(chan struct{})

	processor := GameProcessor{logger, stream, wait, db}

	email := "game-processor-test-2.1@charcoal.sizethree.cc"
	ownerEmail := "game-processor-test-2.2@charcoal.sizethree.cc"
	game, user, owner := models.Game{Status: "ACTIVE", Population: 10}, models.User{Email: email}, models.User{Email: ownerEmail}

	db.Create(&owner)
	db.Create(&user)

	game.OwnerID = owner.ID
	db.Create(&game)

	defer db.Unscoped().Delete(&owner)
	defer db.Unscoped().Delete(&user)
	defer db.Unscoped().Delete(&game)

	stream <- Message{&user, &game, fmt.Sprintf("games:%s", GameProcessorUserLeft)}
	go processor.Begin(ProcessorConfig{DB: testutils.DBConfig()})
	close(stream)
	<-wait

	g := models.Game{}
	db.First(&g, game.ID)

	if g.Population == 9 {
		return
	}

	t.Fatalf("expected game w/ initial population of 10 to decrease to 9 after update")
}
