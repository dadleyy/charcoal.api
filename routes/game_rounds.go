package routes

import "fmt"
import "strconv"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"

func CreateGameRound(runtime *net.RequestRuntime) error {
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

	runtime.AddResult(round.Public())

	return nil
}

func UpdateGameRound(runtime *net.RequestRuntime) error {
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
		return runtime.AddError(e)
	}

	runtime.AddResult(round.Public())

	return nil
}

func DestroyGameRound(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.AddError(fmt.Errorf("BAD_ID"))
	}

	var round models.GameRound

	if e := runtime.First(&round, id).Error; e != nil {
		runtime.Infof("round not found: %d (%s)", id, e.Error())
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	manager, err := runtime.Game(round.GameID)

	if err != nil {
		runtime.Warnf("unable to load game manager: %s", err.Error())
		return runtime.LogicError("invalid-game")
	}

	if manager.IsMember(runtime.User) == false {
		runtime.Infof("user %d not member of game %d", runtime.User.ID, round.GameID)
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	if e := runtime.Delete(&round).Error; e != nil {
		runtime.Infof("failed deletion of round %d: %s", round.ID, e.Error())
		return runtime.AddError(fmt.Errorf("FAILED_DELETE"))
	}

	return nil
}

func FindGameRounds(runtime *net.RequestRuntime) error {
	cursor, results := runtime.Model(&models.GameRound{}), make([]models.GameRound, 0)

	blueprint := runtime.Blueprint(cursor)
	total, err := blueprint.Apply(&results)

	if err != nil {
		runtime.Debugf("invalid blueprint apply: %s", err.Error())
		return err
	}

	for _, r := range results {
		runtime.AddResult(r.Public())
	}

	runtime.SetTotal(total)

	return nil
}
