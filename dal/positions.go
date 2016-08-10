package dal

import "errors"

import "github.com/sizethree/meritoss.api/db"
import "github.com/sizethree/meritoss.api/models"

type PositionFacade struct {
	User uint
	Location int
	Proposal uint
}

func (pos *PositionFacade) IsDuplicate(dbclient *db.Client) bool {
	var existing models.Position
	result := dbclient.Where("user =  ? AND proposal = ?", pos.User, pos.Proposal).First(&existing)
	return result.RecordNotFound() != true
}

func CreatePosition(dbclient *db.Client, facade *PositionFacade) (models.Position, error) {
	var position models.Position
	user, proposal, location := facade.User, facade.Proposal, facade.Location

	if user <  1 {
		return position, errors.New("invalid user id")
	}

	if proposal < 1 {
		return position, errors.New("invalid proposal id")
	}

	if location < -1 || location > 1 {
		return position, errors.New("invalid location")
	}

	var prop models.Proposal
	if e := dbclient.Where("id = ?", facade.Proposal).First(&prop).Error; e != nil {
		return position, errors.New("proposal not found")
	}

	if facade.IsDuplicate(dbclient) {
		return position, errors.New("duplicate position, need to update")
	}

	position = models.Position{Location: location, User: user, Proposal: proposal}
	if e := dbclient.Save(&position).Error; e != nil {
		return position, e
	}

	return position, nil
}
