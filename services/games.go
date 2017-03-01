package services

import "fmt"
import "strconv"
import "net/url"
import "github.com/jinzhu/gorm"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/activity"

const GameManagerInvalidPresident = "INVALID_PRESIDENT"
const GameManagerInvalidVicePresident = "INVALID_VICE_PRESIDENT"
const GameManagerInvalidAsshole = "INVALID_ASSHOLE"

type GameManager struct {
	*gorm.DB
	*log.Logger

	Sockets <-chan activity.Message
	Game    models.Game
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

		if len(value[0]) >= 1 && value[0] != "null" && err != nil {
			m.Infof("invalid value for %s: %v (%v)", rank, value, err)
			return fmt.Errorf("invalid value for %s: %v", rank, value)
		}

		m.Debugf("will be performing update on %d: %s -> [%v]", round.ID, rank, value[0])

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

	return nil
}

func (m *GameManager) IsMember(user models.User) bool {
	member, game := models.GameMembership{}, m.Game

	if e := m.Where("user_id = ? and game_id = ?", user.ID, game.ID).First(&member).Error; e != nil {
		return false
	}

	return true
}

func (m *GameManager) AddUser(user models.User) error {
	member := models.GameMembership{UserID: user.ID, GameID: m.Game.ID}

	if m.IsMember(user) {
		return fmt.Errorf("already a member of the game")
	}

	return m.Create(&member).Error
}
