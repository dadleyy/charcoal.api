package routes

import "fmt"
import "strconv"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"

func CreateGameRound(runtime *net.RequestRuntime) error {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_REQUEST"))
	}

	id, err := strconv.Atoi(body.Get("game"))

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

	round, updates := models.GameRound{}, make(map[string]interface{})

	if err := runtime.First(&round, id).Error; err != nil {
		runtime.Debugf("unable to find round[%d]: %s", id, err.Error())
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	if body.KeyExists("president_id") {
		president, err := strconv.Atoi(body.Get("president_id"))

		if err != nil {
			return runtime.AddError(fmt.Errorf("INVALID_PRESIDENT"))
		}

		member, cursor := models.GameMembership{}, runtime.Where("user_id = ? AND game_id = ?", president, round.GameID)

		if err := cursor.First(&member).Error; err != nil {
			runtime.Debugf("user[%d] not in game: %d, cannot be pres | %s", president, round.GameID, err.Error())
			return runtime.AddError(fmt.Errorf("INVALID_PRESIDENT"))
		}

		updates["president_id"] = president
	}

	if err := runtime.Model(&round).Updates(updates).Error; err != nil {
		runtime.Debugf("failed updating round: %s", err.Error())
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
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
