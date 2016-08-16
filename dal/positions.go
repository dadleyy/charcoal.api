package dal

import "errors"
import "github.com/golang/glog"

import "github.com/sizethree/meritoss.api/db"
import "github.com/sizethree/meritoss.api/models"

type PositionFacade struct {
	User uint
	Location int
	Proposal uint
	ID uint
}

func (pos *PositionFacade) IsDuplicate(dbclient *db.Client) bool {
	var existing models.Position
	result := dbclient.Where("user =  ? AND proposal = ?", pos.User, pos.Proposal).First(&existing)
	return result.RecordNotFound() != true
}

func (pos *PositionFacade) ValidLocation() bool {
	return (pos.Location <= 1) && (pos.Location >= -1)
}

func UpdatePosition(dbclient *db.Client, facade *PositionFacade) error {
	user, id, location := facade.User, facade.ID, facade.Location

	// make sure we have a real position
	var position models.Position
	if e := dbclient.Where("id = ?", id).First(&position).Error; e != nil {
		return nil
	}

	// make sure the position matches that of the user doing the update
	if position.User != user {
		return errors.New("not allowed")
	}

	if !facade.ValidLocation() {
		return errors.New("invalid location")
	}

	// save the position's current location
	previous := position.Location

	// if the last location is the same as this one, avoid doing anything
	if previous == location {
		return nil
	}

	// update the positon's location
	position.Location = location
	if e := dbclient.Save(&position).Error; e != nil {
		return e
	}

	glog.Infof("position changed successfully, creating position history item")
	history := models.PositionHistory{
		Before: previous,
		After: location,
		Proposal: position.Proposal,
		User: position.User,
	}

	return dbclient.Save(&history).Error
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

	if facade.ValidLocation() == false {
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
