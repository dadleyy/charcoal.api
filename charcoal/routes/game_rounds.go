package routes

import "fmt"
import "strconv"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/charcoal/net"
import "github.com/dadleyy/charcoal.api/charcoal/models"

func CreateGameRound(runtime *net.RequestRuntime) *net.ResponseBucket {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		runtime.Warnf("failed body parse: %s", err.Error())
		return runtime.LogicError("invalid-request")
	}

	id, err := strconv.Atoi(body.Get("game_id"))

	if err != nil {
		return runtime.LogicError("invalid-game-id")
	}

	var game models.Game

	if err := runtime.First(&game, id).Error; err != nil {
		runtime.Debugf("unable to find game %d: %s", id, err.Error())
		return runtime.LogicError("game-not-found")
	}

	members, cursor := []models.GameMembership{}, runtime.Where("user_id = ? and game_id = ?", runtime.User.ID, id)

	if err := cursor.Find(&members).Error; err != nil || len(members) == 0 {
		runtime.Debugf("user[%d] does not belong to game[%d]", runtime.User.ID, id)

		if err != nil {
			runtime.Debugf("error: ", err.Error())
		}

		return runtime.LogicError("game-not-found")
	}

	round := models.GameRound{Game: game}

	if err := runtime.Create(&round).Error; err != nil {
		runtime.Debugf("failed saving new round: %s", err.Error())
		return runtime.ServerError()
	}

	return runtime.SendResults(1, []models.GameRound{round})
}

func UpdateGameRound(runtime *net.RequestRuntime) *net.ResponseBucket {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.LogicError("bad-round")
	}

	body, err := forms.Parse(runtime.Request)

	if err != nil {
		runtime.Debugf("error parsing round-update body: %s", err.Error())
		return runtime.LogicError("bad-body")
	}

	round := models.GameRound{}

	if err := runtime.First(&round, id).Error; err != nil {
		runtime.Debugf("unable to find round[%d]: %s", id, err.Error())
		return runtime.LogicError("not-found")
	}

	manager, err := runtime.Game(round.GameID)

	if err != nil {
		runtime.Debugf("unable to find game[%d]: %s", round.GameID, err.Error())
		return runtime.LogicError("not-found")
	}

	if manager.IsMember(runtime.User) == false && runtime.IsAdmin() == false {
		runtime.Debugf("user %d is not in game %d, cannot update", runtime.User.ID, manager.Game.ID)
		return runtime.LogicError("not-found")
	}

	if e := manager.UpdateRound(&round, body.Values); e != nil {
		runtime.Errorf("[update game round] unable to update: %s", e.Error())
		return runtime.ServerError()
	}

	return runtime.SendResults(1, []models.GameRound{round})
}

func DestroyGameRound(runtime *net.RequestRuntime) *net.ResponseBucket {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.SendErrors(fmt.Errorf("bad-id"))
	}

	var round models.GameRound

	if e := runtime.First(&round, id).Error; e != nil {
		runtime.Infof("round not found: %d (%s)", id, e.Error())
		return runtime.SendErrors(fmt.Errorf("not-found"))
	}

	manager, err := runtime.Game(round.GameID)

	if err != nil {
		runtime.Warnf("unable to load game manager: %s", err.Error())
		return runtime.LogicError("invalid-game")
	}

	if manager.IsMember(runtime.User) == false {
		runtime.Infof("user %d not member of game %d", runtime.User.ID, round.GameID)
		return runtime.SendErrors(fmt.Errorf("not-found"))
	}

	if e := runtime.Delete(&round).Error; e != nil {
		runtime.Errorf("[delete game round] failed deletion of round %d: %s", round.ID, e.Error())
		return runtime.ServerError()
	}

	return nil
}

func FindGameRounds(runtime *net.RequestRuntime) *net.ResponseBucket {
	cursor, results := runtime.Model(&models.GameRound{}), make([]models.GameRound, 0)

	blueprint := runtime.Blueprint(cursor)
	total, err := blueprint.Apply(&results)

	if err != nil {
		runtime.Warnf("[find rounds] invalid blueprint apply: %s", err.Error())
		return runtime.LogicError("bad-request")
	}

	return runtime.SendResults(total, results)
}
