package activity

import "fmt"
import "sync"

import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/sizethree/miritos.api/server"
import "github.com/sizethree/miritos.api/models"

type Processor struct {
	Queue chan Message
	DatabaseInfo server.DatabaseConfig
}

func create(message Message, conf *server.DatabaseConfig, out chan<- error) {
	db, err := gorm.Open("mysql", conf.String())

	if err != nil {
		out <- err
		return
	}

	item := models.Activity{
		Type: message.Verb,
		ActorUrl: message.Actor.Url(),
		ActorType: message.Actor.Type(),
		ObjectUrl: message.Object.Url(),
		ObjectType: message.Object.Type(),
	}

	if err := db.Create(&item).Error; err != nil {
		out <- err
	}

	out <- nil
}

func (engine *Processor) Begin() {
	var deferred sync.WaitGroup
	makers := make(chan error)

	for message := range engine.Queue {
		deferred.Add(1)
		go create(message, &engine.DatabaseInfo, makers)
	}

	go func() {
		deferred.Wait()
		close(makers)
	}()

	for err := range makers {
		if err == nil { continue }
		fmt.Errorf("ERROR: %s", err.Error())
	}
}
