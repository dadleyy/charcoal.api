package activity

import "fmt"
import "sync"

import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/mysql"

import "github.com/sizethree/miritos.api/server"
import "github.com/sizethree/miritos.api/models"

type Processor struct {
	Queue chan Message
}

func create(message Message, db *server.Database, out chan<- error) {
	item := models.Activity{
		Type: message.Verb,
		ActorUrl: message.Actor.Url(),
		ActorType: message.Actor.Type(),
		ObjectUrl: message.Object.Url(),
		ObjectType: message.Object.Type(),
	}

	if err := db.Create(&item).Error; err != nil {
		out <- err
		return
	}

	out <- nil
}

func (engine *Processor) Begin(dbconf server.DatabaseConfig) {
	var deferred sync.WaitGroup
	makers := make(chan error)
	conn, err := gorm.Open("mysql", dbconf.String())

	if err != nil {
		panic(fmt.Errorf("BAD_DB_CONFIG"))
	}

	db := server.Database{conn}

	defer db.Close()

	for message := range engine.Queue {
		deferred.Add(1)
		go create(message, &db, makers)
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
