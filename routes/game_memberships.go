package routes

import "fmt"
import "strconv"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"

func DestroyGameMembership(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.AddError(fmt.Errorf("BAD_ID"))
	}

	var membership models.GameMembership

	if err := runtime.First(&membership, id).Error; err != nil {
		runtime.Debugf("error looking for membership: %s", err.Error())
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	if runtime.IsAdmin() == false && membership.UserID != runtime.User.ID {
		runtime.Debugf("cannot delete membership - user[%d] isn't owner", runtime.User.ID)
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	if err := runtime.Delete(&membership).Error; err != nil {
		runtime.Debugf("problem deleting membership: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED_DELETE"))
	}

	return nil
}

func CreateGameMembership(runtime *net.RequestRuntime) error {
	body, err := runtime.Form()

	if err != nil {
		runtime.Debugf("error parsing round-update body: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_REQUEST"))
	}

	game, user := models.Game{}, models.User{}

	if id, err := strconv.Atoi(body.Get("game_id")); err != nil || runtime.First(&game, id).Error != nil {
		runtime.Debugf("invalid game id: %v", body.Get("game_id"))
		return runtime.AddError(fmt.Errorf("INVALID_GAME"))
	}

	if id, err := strconv.Atoi(body.Get("user_id")); err != nil || runtime.First(&user, id).Error != nil {
		runtime.Debugf("invalid user id: %v", body.Get("user_id"))
		return runtime.AddError(fmt.Errorf("INVALID_USER"))
	}

	manager := services.GameManager{runtime.DB, game}

	if err := manager.AddUser(user); err != nil {
		runtime.Debugf("failed adding user: %s", err.Error())
		return runtime.AddError(fmt.Errorf("FAILED_ADD"))
	}

	return nil
}

func FindGameMemberships(runtime *net.RequestRuntime) error {
	cursor, results := runtime.DB, make([]models.GameMembership, 0)
	blueprint := runtime.Blueprint(cursor)

	total, err := blueprint.Apply(&results)

	if err != nil {
		runtime.Debugf("invalid blueprint apply: %s", err.Error())
		return err
	}

	for _, r := range results {
		runtime.AddResult(r)
	}

	runtime.SetTotal(total)

	return nil
}
