package routes

import "fmt"
import "github.com/albrow/forms"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"

func CreateGame(runtime *net.RequestRuntime) error {
	_, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.AddError(err)
	}

	game := models.Game{Owner: runtime.User, Status: models.GameDefaultStatus}

	if err := runtime.Create(&game).Error; err != nil {
		runtime.Errorf("failed saving new game: %s", err.Error())
		return runtime.AddError(err)
	}

	membership := models.GameMembership{User: runtime.User, Game: game}

	if err := runtime.Create(&membership).Error; err != nil {
		runtime.Errorf("unable to create initial membership: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_GAME_CREATE"))
	}

	runtime.AddResult(game)

	return nil
}

func DestroyGame(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.AddError(fmt.Errorf("BAD_ID"))
	}

	var game models.Game

	if err := runtime.First(&game, id).Error; err != nil {
		runtime.Debugf("error looking for game: %s", err.Error())
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	if runtime.IsAdmin() == false && game.OwnerID != runtime.User.ID {
		runtime.Debugf("cannot delete game - user[%d] isn't owner", runtime.User.ID)
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	if err := runtime.Model(&game).Update("status", "ENDED").Error; err != nil {
		runtime.Debugf("problem deleting game: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED_DELETE"))
	}

	return nil
}

func FindGames(runtime *net.RequestRuntime) error {
	var results []models.Game
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&results)

	if err != nil {
		fmt.Errorf("failed game lookup: %s", err.Error())
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	for _, item := range results {
		runtime.AddResult(item)
	}

	runtime.SetTotal(total)

	return nil
}
