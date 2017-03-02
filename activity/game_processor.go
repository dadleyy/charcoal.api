package activity

import "strings"
import "github.com/jinzhu/gorm"
import "github.com/labstack/gommon/log"
import "github.com/dadleyy/charcoal.api/models"

const GameProcessorVerbPrefix = "games:"
const GameProcessorUserJoined = "user-joined"
const GameProcessorUserLeft = "user-left"

type GameProcessor struct {
	*log.Logger
	Stream  chan Message
	Exhaust chan struct{}

	db *gorm.DB
}

func (engine *GameProcessor) updatePopulation(msg Message, amt int) {
	_, vu := msg.Actor.(*models.User)
	game, vg := msg.Object.(*models.Game)

	if vu != true || vg != true {
		engine.Warnf("game user event with invalid object/actor")
		return
	}

	if e := engine.db.Model(&game).Update("population", game.Population+amt).Error; e != nil {
		engine.Errorf("unable to update game population: %s", e.Error())
	}
}

func (engine *GameProcessor) Begin(config ProcessorConfig) {
	database, err := gorm.Open("mysql", config.DB.String())

	if engine.Stream == nil {
		engine.Warnf("no channel provided to the game processor")
		return
	}

	if err != nil {
		panic(err)
	}

	engine.db = database
	defer engine.db.Close()

	for message := range engine.Stream {
		engine.Debugf("received message: %s", message.Verb)
		event := strings.TrimPrefix(message.Verb, GameProcessorVerbPrefix)

		switch event {
		case GameProcessorUserJoined:
			engine.updatePopulation(message, 1)
		case GameProcessorUserLeft:
			engine.updatePopulation(message, -1)
		}

		if engine.Exhaust == nil {
			continue
		}

		engine.Exhaust <- struct{}{}
	}
}
