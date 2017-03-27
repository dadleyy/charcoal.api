package activity

import "fmt"
import "sync"
import "testing"

import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/defs"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/testutils"

func Test_Activity_GameStatsProcessor_UpdateStats(t *testing.T) {
	db, logger, stream := testutils.NewDB(), log.New("miritos"), make(chan Message, 1)

	processor := GameStatsProcessor{
		Logger: logger,
		DB:     db,
		Stream: stream,
	}

	game, users := models.Game{Status: "ACTIVE"}, []models.User{
		models.User{Email: "stat-processor-test-1@sizethree.cc"},
		models.User{Email: "stat-processor-test-2@sizethree.cc"},
		models.User{Email: "stat-processor-test-3@sizethree.cc"},
		models.User{Email: "stat-processor-test-4@sizethree.cc"},
	}

	memberships := []models.GameMembership{}

	for i, _ := range users {
		db.Create(&users[i])
		defer db.Unscoped().Delete(&users[i])
	}

	game.OwnerID = users[0].ID
	db.Create(&game)
	defer db.Unscoped().Delete(&game)

	for _, u := range users {
		m := models.GameMembership{
			Status: "ACTIVE",
			GameID: game.ID,
			UserID: u.ID,
		}

		db.Create(&m)
		defer db.Unscoped().Delete(&m)

		memberships = append(memberships, m)
	}

	ids := []int64{
		int64(users[0].ID),
		int64(users[1].ID),
		int64(users[2].ID),
		int64(users[3].ID),
	}

	rounds := []models.GameRound{
		models.GameRound{
			GameID:          game.ID,
			AssholeID:       &ids[0],
			PresidentID:     &ids[1],
			VicePresidentID: &ids[2],
		},
		models.GameRound{
			GameID:          game.ID,
			AssholeID:       &ids[0],
			PresidentID:     &ids[1],
			VicePresidentID: &ids[3],
		},
		models.GameRound{
			GameID:          game.ID,
			AssholeID:       &ids[0],
			PresidentID:     &ids[1],
			VicePresidentID: &ids[3],
		},
		models.GameRound{
			GameID:          game.ID,
			AssholeID:       &ids[0],
			PresidentID:     &ids[1],
			VicePresidentID: &ids[2],
		},
		models.GameRound{
			GameID:          game.ID,
			AssholeID:       &ids[1],
			PresidentID:     &ids[2],
			VicePresidentID: &ids[3],
		},
		models.GameRound{
			GameID:          game.ID,
			AssholeID:       &ids[0],
			PresidentID:     &ids[1],
			VicePresidentID: &ids[2],
		},
	}

	for i, _ := range rounds {
		db.Create(&rounds[i])
		defer db.Unscoped().Delete(&rounds[i])
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go processor.Begin(&wg)

	stream <- Message{
		Verb:   fmt.Sprintf("%s:%s", defs.GamesStatsStreamIdentifier, defs.GameStatsRoundUpdate),
		Object: &rounds[0],
	}

	close(stream)

	wg.Wait()

	if e := db.Where("game_id = ?", game.ID).Find(&memberships).Error; e != nil {
		t.Fatalf(e.Error())
		return
	}

	if memberships[0].Assholeships != 5 {
		t.Fatalf("expected 4 assholeships but found: %d", memberships[0].Assholeships)
	}

	if memberships[1].Presidencies != 5 {
		t.Fatalf("expected 4 assholeships but found: %d", memberships[0].Assholeships)
	}

	if memberships[2].VicePresidencies != 3 {
		t.Fatalf("expected 4 assholeships but found: %d", memberships[0].Assholeships)
	}
}
