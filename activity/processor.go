package activity

import "fmt"
import "sync"

import "github.com/labstack/gommon/log"

import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/models"

const statusPending string = "PENDING"
const statusApproved string = "APPROVED"
const statusRejected string = "REJECTED"

type Processor struct {
	Queue chan Message
	Logger *log.Logger
	DBConfig db.Config
}

func create(message Message, conn *db.Connection, out chan<- error) {
	item := models.Activity{
		Type: message.Verb,
		ActorUrl: message.Actor.Url(),
		ActorType: message.Actor.Type(),
		ObjectUrl: message.Object.Url(),
		ObjectType: message.Object.Type(),
	}

	if err := conn.Create(&item).Error; err != nil {
		out <- err
		return
	}

	schedule := models.DisplaySchedule{
		Activity: item.ID,
		Approval: statusPending,
	}

	if err := conn.Create(&schedule).Error; err != nil {
		out <- err
		return
	}

	out <- nil
}

func (engine *Processor) Begin() {
	var deferred sync.WaitGroup
	makers := make(chan error)
	conn, err := db.Open(engine.DBConfig)

	if err != nil {
		panic(fmt.Errorf("BAD_DB_CONFIG"))
	}

	defer conn.Close()

	for message := range engine.Queue {
		deferred.Add(1)
		engine.Logger.Debugf("spawning creator goroutine for message: %s", message.Verb)
		go create(message, conn, makers)
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
