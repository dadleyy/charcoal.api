package services

import "fmt"
import "github.com/jinzhu/gorm"
import "github.com/dadleyy/charcoal.api/models"

const GameManagerInvalidPresident = "INVALID_PRESIDENT"
const GameManagerInvalidVicePresident = "INVALID_VICE_PRESIDENT"
const GameManagerInvalidAsshole = "INVALID_ASSHOLE"

type GameManager struct {
	*gorm.DB
	Game models.Game
}

func (m *GameManager) updateRoundRanking(id int, round models.GameRound, ranking string) error {
	user := models.User{}

	if err := m.First(&user, id).Error; err != nil {
		return fmt.Errorf("invalid %s: user not found", ranking)
	}

	if m.IsMember(user) == false {
		return fmt.Errorf("invalid %s: not present in game", ranking)
	}

	switch ranking {
	case "asshole":
		round.AssholeID.Scan(id)
	case "vice_president":
		round.VicePresidentID.Scan(id)
	case "president":
		round.PresidentID.Scan(id)
	default:
		return fmt.Errorf("invalid ranking")
	}

	if err := m.Save(&round).Error; err != nil {
		return fmt.Errorf("unable to save: %s", err.Error())
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

func (m *GameManager) UpdateAsshole(id int, round models.GameRound) error {
	return m.updateRoundRanking(id, round, "asshole")
}

func (m *GameManager) UpdateVicePresident(id int, round models.GameRound) error {
	return m.updateRoundRanking(id, round, "vice_president")
}

func (m *GameManager) UpdatePresident(id int, round models.GameRound) error {
	return m.updateRoundRanking(id, round, "president")
}
