package activity

import "fmt"
import "sync"

import "github.com/labstack/gommon/log"

import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/models"

const statusPending string = "PENDING"
const statusApproved string = "APPROVED"
const statusRejected string = "REJECTED"

type semaphore chan struct{}

type Processor struct {
	Queue chan Message
	Logger *log.Logger
	DBConfig db.Config
}

func create(message Message, conn *db.Connection, out chan<- error, waitlist semaphore) {
	// wait until someone has released a token from the semaphore
	waitlist <- struct{}{}

	item := models.Activity{
		Type: message.Verb,
		ActorUrl: message.Actor.Url(),
		ActorType: message.Actor.Type(),
		ObjectUrl: message.Object.Url(),
		ObjectType: message.Object.Type(),
	}

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

func (engine *Processor) Begin() {
	var deferred sync.WaitGroup
	makers := make(chan error)
	waitlist := make(semaphore, 20)
	conn, err := db.Open(engine.DBConfig)

	if err != nil {
		panic(fmt.Errorf("BAD_DB_CONFIG"))
	}

	defer conn.Close()

	for message := range engine.Queue {
		deferred.Add(1)
		engine.Logger.Debugf("spawning creator goroutine for message: %s", message.Verb)
		go create(message, conn, makers, waitlist)
	}

	go func() {
		deferred.Wait()
		close(makers)
	}()

	for err := range makers {
		if err == nil { continue }
		engine.Logger.Errorf("ERROR: %s", err.Error())
	}
}
