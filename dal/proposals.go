package dal

import "errors"
import "github.com/golang/glog"

import "github.com/sizethree/meritoss.api/db"
import "github.com/sizethree/meritoss.api/models"
import "github.com/sizethree/meritoss.api/middleware"

type ProposalFacade struct {
	Summary string
	Content string
	Author uint
}

func (prop *ProposalFacade) Error() error {
	if len(prop.Summary) < 1 {
		return errors.New("proposal summaries cannot be empty")
	}

	if len(prop.Content) < 1 {
		return errors.New("proposal content fields cannot be empty")
	}

	return nil
}

// FindProposals
// 
// given a database client and a blueprint, returns the list of appro
func FindProposals(client *db.Client, blueprint* middleware.Blueprint) ([]models.Proposal, int, error) {
	var proposals []models.Proposal

	total, e := blueprint.Apply(&proposals, client)

	if e != nil {
		glog.Errorf("errror applying proposal blueprint %s\n", e.Error())
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

	if e := facade.Error(); e != nil {
		return proposal, facade.Error()
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
