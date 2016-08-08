package dal

import "errors"
import "github.com/golang/glog"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/db"
import "github.com/sizethree/meritoss.api/api/models"

type ProposalFacade struct {
	Summary string
	Content string
	Author uint
}

// FindProposals
// 
// given a database client and a blueprint, returns the list of appro
func FindProposals(client *db.Client, blueprint* api.Blueprint) ([]models.Proposal, int, error) {
	var proposals []models.Proposal

	total, e := blueprint.Apply(&proposals, client)

	if e != nil {
		glog.Errorf("errror %s\n", e.Error())
		return proposals, -1, e
	}

	return proposals, total, nil
}


// CreateProposal
//
// given a client and a proposal facade, attempts to create a single proposal
func CreateProposal(client *db.Client, facade *ProposalFacade) (models.Proposal, error) {
	var author models.User
	var proposal models.Proposal

	head := client.Where("ID = ?", facade.Author).Find(&author)

	if e := head.Error; e != nil {
		glog.Errorf("error when finding author %d, %s\n", facade.Author, e.Error())
		return proposal, errors.New("bad author")
	}

	if head.RecordNotFound() {
		glog.Errorf("missing or error when finding author %d\n", facade.Author)
		return proposal, errors.New("bad author")
	}

	proposal = models.Proposal{
		Author: author.ID,
		Summary: facade.Summary,
		Content: facade.Content,
	}

	if err := client.Set("gorm:save_associations", false).Save(&proposal).Error; err != nil {
		glog.Errorf("unable to associate proposal with author: %s\n", err.Error())
		return proposal, err
	}

	glog.Infof("successfully created new proposal #%d\n", proposal.ID)
	return proposal, nil
}

