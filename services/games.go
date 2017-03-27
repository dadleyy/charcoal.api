package services

import "fmt"
import "strconv"
import "net/url"
import "github.com/jinzhu/gorm"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/defs"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/activity"

type GameManager struct {
	*gorm.DB
	*log.Logger

	Streams map[string](chan<- activity.Message)
	Game    models.Game
}

func (m *GameManager) OwnerID() uint {
	return m.Game.OwnerID
}

func (m *GameManager) IsEnded() bool {
	return m.Game.Status == "ENDED" || m.Game.DeletedAt != nil
}

func (m *GameManager) EndGame() error {
	if m.IsEnded() {
		m.Infof("game already ended, skipping.")
		return nil
	}

	e := m.Model(&m.Game).Update("status", "ENDED").Error

	if e != nil {
		return e
	}

	stream, exists := m.Streams[defs.GamesStreamIdentifier]

	if exists != true {
		return nil
	}

	owner := models.User{}

	if e := m.First(&owner, m.Game.OwnerID).Error; e != nil {
		m.Warnf("unable to publish ended event due to owner lookup: %s", e.Error())
		return nil
	}

	verb := fmt.Sprintf("%s:%s", defs.GamesStreamIdentifier, defs.GameProcessorGameEnded)
	stream <- activity.Message{&owner, &m.Game, verb}

	return nil
}

func (m *GameManager) ApplyUpdates(updates url.Values) error {
	fields := make(map[string]interface{})

	if status := updates.Get("status"); status == "ENDED" {
		return m.EndGame()
	}

	if name := updates.Get("name"); len(name) >= 2 {
		fields["name"] = name
	}

	m.Debugf("applying updates: %v", fields)
	return m.Model(&m.Game).Update(fields).Error
}

func (m *GameManager) UpdateRound(round *models.GameRound, rankings url.Values) error {
	updates := make(map[string]*int64)

	if round == nil || round.ID >= 1 == false {
		return fmt.Errorf("invalid round id")
	}

	for _, rank := range []string{"vice_president_id", "president_id", "asshole_id"} {
		value, exists := rankings[rank]

		if exists == false || len(value) != 1 {
			continue
		}

		id, err := strconv.ParseInt(value[0], 10, 64)

		// If we have a value that is not null and there is an error parsing the int, we have an invalid value
		if len(value[0]) >= 1 && value[0] != "null" && err != nil {
			m.Infof("[game processor] invalid value for %s: %v (%v)", rank, value, err)
			return fmt.Errorf("invalid value for %s: %v", rank, value)
		}

		m.Debugf("[game processor] will be performing update on %d: %s -> [%v]", round.ID, rank, value[0])

		// Clearing out the value
		if err != nil {
			updates[rank] = nil
			continue
		}

		u := models.User{Common: models.Common{ID: uint(id)}}

		if valid := m.IsMember(u); valid != true {
			return fmt.Errorf("user %d not a member of game %d", u.ID, m.Game.ID)
		}

		updates[rank] = &id
	}

	if e := m.Model(round).Update(updates).Error; e != nil {
		return e
	}

	if s, ok := m.Streams[defs.GamesStatsStreamIdentifier]; ok == true {
		m.Debugf("[game service] publishing to stats stream")
		verb := fmt.Sprintf("%s:%s", defs.GamesStatsStreamIdentifier, defs.GameStatsRoundUpdate)
		s <- activity.Message{Verb: verb, Object: round}
	}

	return nil
}

func (m *GameManager) IsMember(user models.User) bool {
	member, game := models.GameMembership{}, m.Game
	cursor := m.Where("user_id = ? and game_id = ?", user.ID, game.ID)

	if e := cursor.Where("status = ?", defs.GameMembershipActiveStatus).First(&member).Error; e != nil {
		return false
	}

	return true
}

func (m *GameManager) RemoveMember(member models.GameMembership) error {
	if m.Game.ID != member.GameID {
		return fmt.Errorf("game/membership mismatch")
	}

	user := models.User{}

	if e := m.First(&user, member.UserID).Error; e != nil {
		return e
	}

	if e := m.Model(&member).Update("status", defs.GameMembershipInactiveStatus).Error; e != nil {
		m.Errorf("[game manager] unable to remove member: %s", e.Error())
		return e
	}

	m.Debugf("removed member: %d from game %d", member.UserID, m.Game.ID)

	if stream, ok := m.Streams[defs.GamesStreamIdentifier]; ok {
		verb := fmt.Sprintf("%s:%s", defs.GamesStreamIdentifier, defs.GameProcessorUserLeft)
		stream <- activity.Message{&user, &m.Game, verb}
	}

	return nil
}

func (m *GameManager) AddUser(user models.User) (models.GameMembership, error) {
	member := models.GameMembership{}

	if m.IsMember(user) {
		return models.GameMembership{}, fmt.Errorf("already a member of the game")
	}

	publish := func() {
		stream, ok := m.Streams[defs.GamesStreamIdentifier]
		verb := fmt.Sprintf("%s:%s", defs.GamesStreamIdentifier, defs.GameProcessorUserJoined)

		if ok == true {
			stream <- activity.Message{&user, &m.Game, verb}
		}
	}

	// At this point, no matter what happens we will want to recalculate the population of the game.
	defer publish()

	// If we already have a membership record associated w/ the game and user, just update it so that it reflects an
	// active status and carry on.
	if m.Where("user_id = ? AND game_id = ?", user.ID, m.Game.ID).First(&member).RecordNotFound() != true {
		if e := m.Model(&member).Update("status", defs.GameMembershipActiveStatus).Error; e != nil {
			return models.GameMembership{}, e
		}

		return member, nil
	}

	member = models.GameMembership{
		UserID: user.ID,
		GameID: m.Game.ID,
		Status: defs.GameMembershipActiveStatus,
	}

	if e := m.Create(&member).Error; e != nil {
		return models.GameMembership{}, e
	}

	return member, nil
}
