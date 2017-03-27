package activity

import "fmt"
import "sync"
import "strings"
import "github.com/jinzhu/gorm"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/defs"
import "github.com/dadleyy/charcoal.api/models"

type GameProcessor struct {
	*log.Logger
	*gorm.DB

	Stream chan Message
}

func (engine *GameProcessor) Begin(wg *sync.WaitGroup) {
	defer wg.Done()

	if engine.Stream == nil {
		engine.Warnf("no channel provided to the game processor")
		return
	}

	engine.Debugf("[games processor] starting game processor")
	internal := sync.WaitGroup{}

	for message := range engine.Stream {
		engine.Debugf("[game bg processor] received message: %s", message.Verb)
		event := strings.TrimPrefix(message.Verb, fmt.Sprintf("%s:", defs.GamesStreamIdentifier))

		internal.Add(1)

		switch event {
		case defs.GameProcessorUserJoined:
			engine.playerJoined(message, &internal)
		case defs.GameProcessorUserLeft:
			engine.playerLeft(message, &internal)
		default:
			internal.Done()
		}
	}

	internal.Wait()
}

func (engine *GameProcessor) userPayload(msg Message) (*models.User, *models.Game, bool) {
	user, vu := msg.Actor.(*models.User)
	game, vg := msg.Object.(*models.Game)
	return user, game, (vu && vg)
}

func (engine *GameProcessor) playerJoined(msg Message, wg *sync.WaitGroup) {
	defer wg.Done()
	user, game, ok := engine.userPayload(msg)

	if !ok {
		return
	}

	defer engine.updatePopulation(game)

	rounds, history := []models.GameRound{}, models.GameMembershipHistory{GameID: game.ID, UserID: user.ID}

	if e := engine.Where("game_id = ?", game.ID).Limit(1).Order("id DESC").Find(&rounds).Error; e != nil {
		engine.Errorf("[game processor] unable to find rounds: %s", e.Error())
		return
	}

	if len(rounds) >= 1 {
		history.EntryRoundID = &rounds[0].ID
		engine.Debugf("[game processor] late join: game[%s] user[%s] round[%d]", game.Uuid, user.Uuid, rounds[0].ID)
	}

	if e := engine.Create(&history).Error; e != nil {
		engine.Errorf("[game processor] unable to create join history: %s", e.Error())
		return
	}

	engine.Debugf("[game processor] player joined, created history record: %s", history.Uuid)
}

func (engine *GameProcessor) playerLeft(msg Message, wg *sync.WaitGroup) {
	defer wg.Done()
	user, game, ok := engine.userPayload(msg)

	if !ok {
		return
	}

	defer engine.updatePopulation(game)

	rounds, history := []models.GameRound{}, []models.GameMembershipHistory{}
	cursor := engine.Where("game_id = ? AND user_id = ?", game.ID, user.ID)

	if e := cursor.Limit(1).Order("id DESC").Find(&history).Error; e != nil {
		engine.Errorf("[game processor] unable to lookup history records for game, err: %s", e.Error())
		return
	}

	if len(history) >= 1 == false {
		engine.Warnf("[game processor] no history record for user[%d] in game[%d]", user.ID, game.ID)
		return
	}

	if e := engine.Where("game_id = ?", game.ID).Limit(1).Order("id DESC").Find(&rounds).Error; e != nil {
		engine.Errorf("[game processor] unable to find rounds: %s", e.Error())
		return
	}

	if len(rounds) >= 1 == false {
		engine.Warnf("[game processor] user left before game began, deleting useless history")
		engine.Unscoped().Delete(&history)
		return
	}

	if e := engine.Model(&history[0]).Update("exit_round_id", rounds[0].ID).Error; e != nil {
		engine.Errorf("[game processor] unable to update history record: %s", e.Error())
		return
	}

	engine.Debugf("[game processor] player left, updated history record: %s", history[0].Uuid)
}

func (engine *GameProcessor) updatePopulation(game *models.Game) {
	count, membership := 0, models.GameMembership{}
	cursor := engine.Model(&membership).Where("game_id = ? AND status = ?", game.ID, defs.GameMembershipActiveStatus)

	if e := cursor.Count(&count).Error; e != nil {
		engine.Errorf("[game processor] unable to get game population count: %s", e.Error())
		return
	}

	if e := engine.Model(&game).Update("population", count).Error; e != nil {
		engine.Errorf("[game processor] unable to update game population: %s", e.Error())
		return
	}

	engine.Debugf("[game processor] population updated!")
}
