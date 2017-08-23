package routes

import "strconv"

import "github.com/dadleyy/charcoal.api/charcoal/net"
import "github.com/dadleyy/charcoal.api/charcoal/defs"
import "github.com/dadleyy/charcoal.api/charcoal/models"

// DestroyGameMembership deletes a membership record if the current user is able to do so.
func DestroyGameMembership(runtime *net.RequestRuntime) *net.ResponseBucket {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.LogicError("invalid-id")
	}

	var membership models.GameMembership

	if err := runtime.First(&membership, id).Error; err != nil {
		runtime.Debugf("[delete membership] error looking for membership: %s", err.Error())
		return runtime.LogicError("not-found")
	}

	manager, err := runtime.Game(membership.GameID)

	if err != nil {
		runtime.Warnf("[delete membership] unable to load game manager: %s", err.Error())
		return runtime.LogicError("not-found")
	}

	// membership, current, owner
	m, c, o := membership.UserID, runtime.User.ID, manager.Game.OwnerID

	if runtime.IsAdmin() == false && m != c && c != o {
		runtime.Warnf("[delete membership] cannot delete membership - user[%d] isn't owner", runtime.User.ID)
		return runtime.LogicError("bad-user")
	}

	if e := manager.RemoveMember(membership); e != nil {
		runtime.Errorf("[delete membership] problem deleting membership: %s", e.Error())
		return runtime.ServerError()
	}

	if e := runtime.Delete(&membership).Error; e != nil {
		runtime.Errorf("[delete membership] problem deleting membership: %s", e.Error())
		return runtime.ServerError()
	}

	return nil
}

func CreateGameMembership(runtime *net.RequestRuntime) *net.ResponseBucket {
	body, err := runtime.Form()

	if err != nil {
		runtime.Warnf("[create game mem] error parsing game memebership body: %s", err.Error())
		return runtime.ServerError()
	}

	id, err := strconv.Atoi(body.Get("game_id"))

	if err != nil {
		runtime.Warnf("[create game mem] invalid game id: %v", body.Get("game_id"))
		return runtime.LogicError("bad-game")
	}

	manager, err := runtime.Game(uint(id))

	if err != nil {
		runtime.Warnf("[create game mem] invalid game id: %v", err.Error())
		return runtime.LogicError("bad-game")
	}

	user := models.User{}

	if id, err := strconv.Atoi(body.Get("user_id")); err != nil || runtime.First(&user, id).Error != nil {
		runtime.Warnf("[create game mem] invalid user id: %v", body.Get("user_id"))
		return runtime.LogicError("bad-user")
	}

	if manager.IsMember(runtime.User) == false && runtime.IsAdmin() == false {
		runtime.Warnf("[create game mem] invalid user id: %v", body.Get("user_id"))
		return runtime.LogicError("not-found")
	}

	membership, err := manager.AddUser(user)

	if err != nil {
		runtime.Errorf("[create game mem] failed adding user: %s", err.Error())
		return runtime.LogicError("invalid-membership")
	}

	return runtime.SendResults(1, membership)
}

func UpdateGameMembership(runtime *net.RequestRuntime) *net.ResponseBucket {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.LogicError("invalid-membership")
	}

	body, err := runtime.Form()

	if err != nil {
		return runtime.LogicError("bad-request")
	}

	s := body.Get("status")

	if s != defs.GameMembershipActiveStatus && s != defs.GameMembershipInactiveStatus {
		runtime.Warnf("[update game mem] unable to updating membership to: %s", s)
		return runtime.LogicError("invalid-membership-update")
	}

	membership := models.GameMembership{}

	if e := runtime.First(&membership, id).Error; e != nil {
		runtime.Errorf("[update game mem] unable to find membership: %s", e.Error())
		return runtime.LogicError("bad-membership")
	}

	manager, e := runtime.Game(membership.GameID)

	if e != nil {
		runtime.Errorf("[update game mem] unable to find game: %s", e.Error())
		return runtime.LogicError("bad-game")
	}

	if manager.IsMember(runtime.User) == false && runtime.IsAdmin() == false {
		runtime.Warnf("[update game mem] failed attempt to update membership due to privileges")
		return runtime.LogicError("not-found")
	}

	if s == defs.GameMembershipActiveStatus {
		_, e = manager.AddUser(models.User{Common: models.Common{ID: membership.UserID}})
	} else {
		e = manager.RemoveMember(membership)
	}

	if e != nil {
		runtime.Errorf("[update game mem] unable to update membership: %s", e.Error())
		return runtime.LogicError("bad-request")
	}

	runtime.Debugf("[update game member] updating %d: %v", id, body.Get("status"))

	return nil
}

func FindGameMemberships(runtime *net.RequestRuntime) *net.ResponseBucket {
	cursor, results := runtime.DB, make([]models.GameMembership, 0)
	blueprint := runtime.Blueprint(cursor)

	total, err := blueprint.Apply(&results)

	if err != nil {
		runtime.Errorf("[membership search] invalid blueprint apply: %s", err.Error())
		return runtime.ServerError()
	}

	return runtime.SendResults(total, results)
}
