package routes

import "fmt"
import "github.com/albrow/forms"
import "github.com/docker/docker/pkg/namesgenerator"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/defs"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/activity"

func UpdateGame(runtime *net.RequestRuntime) *net.ResponseBucket {
	id, ok := runtime.IntParam("id")
	body, err := forms.Parse(runtime.Request)

	if err != nil || ok != true {
		runtime.Infof("invalid body received: %s", err.Error())
		return runtime.LogicError("bad-request")
	}

	manager, err := runtime.Game(uint(id))

	if err != nil {
		runtime.Warnf("unable to get manager for game: %s", err.Error())
		return runtime.LogicError("bad-request")
	}

	if manager.IsMember(runtime.User) == false && runtime.IsAdmin() == false {
		runtime.Debugf("invalid user tried to update game %d: %d", manager.Game.ID, runtime.User.ID)
		return runtime.LogicError("not-found")
	}

	if e := manager.ApplyUpdates(body.Values); e != nil {
		runtime.Warnf("unable to save game updates: %s", e.Error())
		return runtime.ServerError()
	}

	return runtime.SendResults(1, []models.Game{manager.Game})
}

func CreateGame(runtime *net.RequestRuntime) *net.ResponseBucket {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.SendErrors(err)
	}

	name := namesgenerator.GetRandomName(0)

	if n := body.Get("name"); len(n) > 4 && len(n) < 20 {
		name = n
	}

	game := models.Game{Name: name, Owner: runtime.User, Status: models.GameDefaultStatus}

	if err := runtime.Create(&game).Error; err != nil {
		runtime.Errorf("[create game] failed saving new game: %s", err.Error())
		return runtime.ServerError()
	}

	membership := models.GameMembership{
		User:   runtime.User,
		Game:   game,
		Status: defs.GameMembershipActiveStatus,
	}

	if err := runtime.Create(&membership).Error; err != nil {
		runtime.Errorf("[create game] unable to create initial membership: %s", err.Error())
		return runtime.ServerError()
	}

	verb := defs.GameProcessorVerbPrefix + defs.GameProcessorUserJoined
	runtime.Publish(activity.Message{&runtime.User, &game, verb})

	return runtime.SendResults(1, []models.Game{game})
}

func DestroyGame(runtime *net.RequestRuntime) *net.ResponseBucket {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.LogicError("invalid-game")
	}

	manager, err := runtime.Game(uint(id))

	if err != nil {
		runtime.Debugf("error looking for game: %s", err.Error())
		return runtime.LogicError("not-found")
	}

	if runtime.IsAdmin() == false && manager.OwnerID() != runtime.User.ID {
		runtime.Debugf("cannot delete game - user[%d] isn't owner", runtime.User.ID)
		return runtime.LogicError("not-found")
	}

	if err := manager.EndGame(); err != nil {
		runtime.Debugf("problem deleting game: %s", err.Error())
		return runtime.ServerError()
	}

	if err := runtime.Delete(&manager.Game).Error; err != nil {
		runtime.Warnf("unable to delete record: %s", err.Error())
		return runtime.ServerError()
	}

	return nil
}

func FindGames(runtime *net.RequestRuntime) *net.ResponseBucket {
	var results []models.Game
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&results)

	if err != nil {
		fmt.Errorf("[find games] failed game lookup: %s", err.Error())
		return runtime.LogicError("bad-request")
	}

	return runtime.SendResults(total, results)
}
