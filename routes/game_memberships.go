package routes

import "strconv"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"

func DestroyGameMembership(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.LogicError("invalid-id")
	}

	var membership models.GameMembership

	if err := runtime.First(&membership, id).Error; err != nil {
		runtime.Debugf("error looking for membership: %s", err.Error())
		return runtime.LogicError("not-found")
	}

	manager, err := runtime.Game(membership.GameID)

	if err != nil {
		runtime.Warnf("unable to load game manager: %s", err.Error())
		return runtime.LogicError("not-found")
	}

	// membership, current, owner
	m, c, o := membership.UserID, runtime.User.ID, manager.Game.OwnerID

	if runtime.IsAdmin() == false && m != c && c != o {
		runtime.Debugf("cannot delete membership - user[%d] isn't owner", runtime.User.ID)
		return runtime.LogicError("bad-user")
	}

	if e := manager.RemoveMember(membership); e != nil {
		runtime.Debugf("problem deleting membership: %s", e.Error())
		return runtime.ServerError()
	}

	return nil
}

func CreateGameMembership(runtime *net.RequestRuntime) error {
	body, err := runtime.Form()

	if err != nil {
		runtime.Debugf("error parsing game memebership body: %s", err.Error())
		return runtime.ServerError()
	}

	id, err := strconv.Atoi(body.Get("game_id"))

	if err != nil {
		runtime.Debugf("invalid game id: %v", body.Get("game_id"))
		return runtime.LogicError("bad-game")
	}

	manager, err := runtime.Game(uint(id))

	if err != nil {
		runtime.Debugf("invalid game id: %v", err.Error())
		return runtime.LogicError("bad-game")
	}

	user := models.User{}

	if id, err := strconv.Atoi(body.Get("user_id")); err != nil || runtime.First(&user, id).Error != nil {
		runtime.Debugf("invalid user id: %v", body.Get("user_id"))
		return runtime.LogicError("bad-user")
	}

	if manager.IsMember(runtime.User) == false && runtime.IsAdmin() == false {
		runtime.Debugf("invalid user id: %v", body.Get("user_id"))
		return runtime.LogicError("not-found")
	}

	membership, err := manager.AddUser(user)

	if err != nil {
		runtime.Debugf("failed adding user: %s", err.Error())
		return runtime.LogicError("invalid-membership")
	}

	runtime.AddResult(membership)
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
