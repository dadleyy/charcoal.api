package routes

import "fmt"
import "strconv"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"

func CreateGameRound(runtime *net.RequestRuntime) error {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_REQUEST"))
	}

	id, err := strconv.Atoi(body.Get("game_id"))

	if err != nil {
		return runtime.AddError(fmt.Errorf("MISSING_GAME_ID"))
	}

	var game models.Game

	if err := runtime.First(&game, id).Error; err != nil {
		runtime.Debugf("unable to find game %d: %s", id, err.Error())
		return runtime.AddError(fmt.Errorf("GAME_NOT_FOUND"))
	}

	members, cursor := []models.GameMembership{}, runtime.Where("user_id = ? and game_id = ?", runtime.User.ID, id)

	if err := cursor.Find(&members).Error; err != nil || len(members) == 0 {
		runtime.Debugf("user[%d] does not belong to game[%d]", runtime.User.ID, id)

		if err != nil {
			runtime.Debugf("error: ", err.Error())
		}

		return runtime.AddError(fmt.Errorf("GAME_NOT_FOUND"))
	}

	round := models.GameRound{Game: game}

	if err := runtime.Create(&round).Error; err != nil {
		runtime.Debugf("failed saving new round: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED_CREATE"))
	}

	runtime.AddResult(round.Public())

	return nil
}

func UpdateGameRound(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.AddError(fmt.Errorf("BAD_ID"))
	}

	body, err := forms.Parse(runtime.Request)

	if err != nil {
		runtime.Debugf("error parsing round-update body: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_REQUEST"))
	}

	round, game := models.GameRound{}, models.Game{}

	if err := runtime.First(&round, id).Error; err != nil {
		runtime.Debugf("unable to find round[%d]: %s", id, err.Error())
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	if err := runtime.First(&game, round.GameID).Error; err != nil {
		runtime.Debugf("unable to find game[%d]: %s", round.GameID, err.Error())
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	manager := services.GameManager{runtime.DB, game}

	if manager.IsMember(runtime.User) == false && runtime.IsAdmin() == false {
		runtime.Debugf("user %d is not in game %d, cannot update", runtime.User.ID, game.ID)
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	if body.KeyExists("vice_president_id") {
		vp, err := strconv.Atoi(body.Get("vice_president_id"))

		if err != nil {
			return runtime.AddError(fmt.Errorf(services.GameManagerInvalidVicePresident))
		}

		runtime.Debugf("up. vp of game[%d], round[%d] to user[%d]", game.ID, round.ID, vp)

		if err := manager.UpdateVicePresident(vp, round); err != nil {
			runtime.Debugf("failed vp update: %s", err.Error())
			return runtime.AddError(fmt.Errorf(services.GameManagerInvalidVicePresident))
		}

		if e := runtime.First(&round).Error; e != nil {
			return runtime.AddError(e)
		}
	}

	if body.KeyExists("asshole_id") {
		ass, err := strconv.Atoi(body.Get("asshole_id"))

		if err != nil {
			return runtime.AddError(fmt.Errorf(services.GameManagerInvalidPresident))
		}

		runtime.Debugf("up. ass of game[%d], round[%d] to user[%d]", game.ID, round.ID, ass)

		if err := manager.UpdateAsshole(ass, round); err != nil {
			runtime.Debugf("failed ass update: %s", err.Error())
			return runtime.AddError(fmt.Errorf(services.GameManagerInvalidPresident))
		}

		if e := runtime.First(&round).Error; e != nil {
			return runtime.AddError(e)
		}
	}

	if body.KeyExists("president_id") {
		president, err := strconv.Atoi(body.Get("president_id"))

		if err != nil {
			return runtime.AddError(fmt.Errorf(services.GameManagerInvalidPresident))
		}

		runtime.Debugf("up. pressy of game[%d], round[%d] to user[%d]", game.ID, round.ID, president)

		if err := manager.UpdatePresident(president, round); err != nil {
			runtime.Debugf("failed pressy update: %s", err.Error())
			return runtime.AddError(fmt.Errorf(services.GameManagerInvalidPresident))
		}

		if e := runtime.First(&round).Error; e != nil {
			return runtime.AddError(e)
		}
	}

	if e := runtime.First(&round).Error; e != nil {
		runtime.AddError(e)
	}

	runtime.AddResult(round.Public())

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
