package activity

import "sync"

import "github.com/jinzhu/gorm"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/models"

type GameStatsProcessor struct {
	*log.Logger
	*gorm.DB

	Stream chan Message
}

func (engine *GameStatsProcessor) Begin(wg *sync.WaitGroup) {
	defer wg.Done()

	if engine.Stream == nil {
		engine.Warnf("[game stats processor] no channel provided to the game processor")
		return
	}

	internal := sync.WaitGroup{}
	engine.Debugf("[game stats processor] starting stats processor")

	for m := range engine.Stream {
		engine.Debugf("[game stats processor] received message: %s", m.Verb)
		internal.Add(1)
		go engine.process(&m, &internal)
	}

	internal.Wait()
}

func (engine *GameStatsProcessor) process(message *Message, wg *sync.WaitGroup) {
	defer wg.Done()
	round, ok := message.Object.(*models.GameRound)

	if ok != true || round == nil {
		engine.Warnf("[game stats] unable to coerce message object to game round!")
		return
	}

	rounds, members := []models.GameRound{}, []models.GameMembership{}

	if e := engine.Where("game_id = ?", round.GameID).Find(&members).Error; e != nil {
		engine.Errorf("[game stats] failed fetching members: %s", e.Error())
		return
	}

	if e := engine.Where("game_id = ?", round.GameID).Find(&rounds).Error; e != nil {
		engine.Errorf("[game stats] failed fetching rounds: %s", e.Error())
		return
	}

	for _, m := range members {
		user, stats := int64(m.UserID), map[string]int{
			"assholeships":      0,
			"presidencies":      0,
			"vice_presidencies": 0,
		}

		for _, r := range rounds {
			if r.AssholeID != nil && *r.AssholeID == user {
				stats["assholeships"]++
			}

			if r.PresidentID != nil && *r.PresidentID == user {
				stats["presidencies"]++
			}

			if r.VicePresidentID != nil && *r.VicePresidentID == user {
				stats["vice_presidencies"]++
			}
		}

		if e := engine.Model(&m).Update(stats).Error; e != nil {
			engine.Errorf("[game stats] unable to update membership stats")
		}
	}

	engine.Debugf("[game stats] updated stats for round[%d]", round.ID)
}
