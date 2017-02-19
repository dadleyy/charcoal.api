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

	game := models.Game{Owner: runtime.User}

	if err := runtime.Database().Create(&game).Error; err != nil {
		runtime.Errorf("failed saving new game: %s", err.Error())
		return runtime.AddError(err)
	}

	membership := models.GameMembership{User: runtime.User, Game: game}

	if err := runtime.Database().Create(&membership).Error; err != nil {
		runtime.Errorf("unable to create initial membership: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_GAME_CREATE"))
	}

	runtime.AddResult(game)

	return nil
}

func FindGames(runtime *net.RequestRuntime) error {
	var results []models.Game
	cursor := runtime.Cursor(&runtime.User)

	if err := cursor.Related(&results, "Games").Error; err != nil {
		fmt.Errorf("failed game lookup: %s", err.Error())
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	for _, item := range results {
		runtime.AddResult(item)
	}

	return nil
}
