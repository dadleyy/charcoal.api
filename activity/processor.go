package activity

import "sync"
import "github.com/labstack/gommon/log"

import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/models"

const statusPending string = "PENDING"
const statusApproved string = "APPROVED"
const statusRejected string = "REJECTED"
const spawnLog = "spawning goroutine: \"%s\" - \"%s\" by \"%s\""

const VerbCreated = "created"
const VerbDeleted = "deleted"

type semaphore chan struct{}

type Processor struct {
	*log.Logger
	DatabaseConnection *db.Connection
	Queue              chan Message
}

func create(message Message, conn *db.Connection, out chan<- error, waitlist semaphore, def *sync.WaitGroup) {
	// wait until someone has released a token from the semaphore
	waitlist <- struct{}{}

	item := models.Activity{
		Type:       message.Verb,
		ActorUrl:   message.Actor.Url(),
		ActorType:  message.Actor.Type(),
		ObjectUrl:  message.Object.Url(),
		ObjectType: message.Object.Type(),
	}

	defer def.Done()

	if err := conn.Create(&item).Error; err != nil {
		out <- err

		// release our token by receiving from the channel
		<-waitlist
		return
	}

	schedule := models.DisplaySchedule{
		Activity: item.ID,
		Approval: statusPending,
	}

	if err := conn.Create(&schedule).Error; err != nil {
		out <- err
		// release our token by receiving from the channel
		<-waitlist
		return
	}

	// release our token by receiving from the channel
	<-waitlist
	out <- nil
}

func destroy(message Message, conn *db.Connection, out chan<- error, waitlist semaphore, def *sync.WaitGroup) {
	// wait until someone has released a token from the semaphore
	waitlist <- struct{}{}

	finish := func() {
		defer def.Done()
		<-waitlist
	}

	defer finish()

	references := make([]struct{ ID uint }, 0)
	cursor := conn.Table("activity").Select("id").Where("object_url = ?", message.Object.Url())

	if err := cursor.Scan(&references).Error; err != nil {
		out <- err
		return
	}

	ids := make([]uint, len(references))
	for index, activity := range references {
		ids[index] = activity.ID
	}

	if err := conn.Unscoped().Delete(models.DisplaySchedule{}, "activity in (?)", ids).Error; err != nil {
		out <- err
		return
	}

	if err := conn.Unscoped().Delete(models.Activity{}, "id in (?)", ids).Error; err != nil {
		out <- err
		return
	}

	out <- nil
}

func (engine *Processor) Begin() {
	var deferred sync.WaitGroup
	makers := make(chan error)
	waitlist := make(semaphore, 20)

	for message := range engine.Queue {
		deferred.Add(1)
		engine.Debugf(spawnLog, message.Verb, message.Object.Url(), message.Actor.Url())

		if message.Verb == VerbDeleted {
			go destroy(message, engine.DatabaseConnection, makers, waitlist, &deferred)
			continue
		}

		go create(message, engine.DatabaseConnection, makers, waitlist, &deferred)
	}

	finish := func() {
		deferred.Wait()
		close(makers)
	}

	go finish()

	for err := range makers {
		if err == nil {
			continue
		}
		engine.Logger.Errorf("ERROR: %s", err.Error())
	}
}
