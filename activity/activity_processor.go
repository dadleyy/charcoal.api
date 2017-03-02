package activity

import "sync"
import "strings"
import "github.com/jinzhu/gorm"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/models"

const statusPending string = "PENDING"
const statusApproved string = "APPROVED"
const statusRejected string = "REJECTED"
const spawnLog = "spawning goroutine: \"%s\" - \"%s\" by \"%s\""

const VerbCreated = "created"
const VerbDeleted = "deleted"

type ActivityProcessor struct {
	*log.Logger
	Stream chan Message

	db *gorm.DB
}

func (engine *ActivityProcessor) create(message Message, out chan<- error, waitlist semaphore, def *sync.WaitGroup) {
	// wait until someone has released a token from the semaphore
	waitlist <- struct{}{}

	item := models.Activity{
		Type:       message.Verb,
		ActorUrl:   message.Actor.Url(),
		ActorType:  message.Actor.Type(),
		ActorUuid:  message.Actor.Identifier(),
		ObjectUrl:  message.Object.Url(),
		ObjectType: message.Object.Type(),
		ObjectUuid: message.Object.Identifier(),
	}

	defer def.Done()

	if err := engine.db.Create(&item).Error; err != nil {
		out <- err

		// release our token by receiving from the channel
		<-waitlist
		return
	}

	schedule := models.DisplaySchedule{
		Activity: item.ID,
		Approval: statusPending,
	}

	if err := engine.db.Create(&schedule).Error; err != nil {
		out <- err
		// release our token by receiving from the channel
		<-waitlist
		return
	}

	// release our token by receiving from the channel
	<-waitlist
	out <- nil
}

func (engine *ActivityProcessor) destroy(message Message, out chan<- error, waitlist semaphore, def *sync.WaitGroup) {
	// wait until someone has released a token from the semaphore
	waitlist <- struct{}{}

	finish := func() {
		defer def.Done()
		<-waitlist
	}

	defer finish()

	references := make([]struct{ ID uint }, 0)
	cursor := engine.db.Table("activity").Select("id").Where("object_url = ?", message.Object.Url())

	if err := cursor.Scan(&references).Error; err != nil {
		out <- err
		return
	}

	ids := make([]uint, len(references))
	for index, activity := range references {
		ids[index] = activity.ID
	}

	if err := engine.db.Unscoped().Delete(models.DisplaySchedule{}, "activity in (?)", ids).Error; err != nil {
		out <- err
		return
	}

	if err := engine.db.Unscoped().Delete(models.Activity{}, "id in (?)", ids).Error; err != nil {
		out <- err
		return
	}

	out <- nil
}

func (engine *ActivityProcessor) Begin(config ProcessorConfig) {
	var deferred sync.WaitGroup
	makers := make(chan error)
	waitlist := make(semaphore, 20)

	database, err := gorm.Open("mysql", config.DB.String())

	if err != nil {
		panic(err)
	}

	engine.db = database
	defer engine.db.Close()

	listen := func() {
		for err := range makers {
			if err == nil {
				continue
			}

			engine.Logger.Errorf("ERROR: %s", err.Error())
		}
	}

	go listen()

	for message := range engine.Stream {
		identifiers := strings.Split(message.Verb, ":")

		if len(identifiers) != 2 || identifiers[0] != "activity" {
			engine.Debugf("skipping published message, not activity: %s", message.Verb)
			continue
		}

		deferred.Add(1)
		engine.Debugf(spawnLog, message.Verb, message.Object.Identifier(), message.Actor.Identifier())

		if message.Verb == VerbDeleted {
			go engine.destroy(message, makers, waitlist, &deferred)
			continue
		}

		go engine.create(message, makers, waitlist, &deferred)
	}

	deferred.Wait()
	close(makers)
	close(waitlist)
}
