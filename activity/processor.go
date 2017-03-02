package activity

import "github.com/dadleyy/charcoal.api/db"

type semaphore chan struct{}

type ProcessorConfig struct {
	DB db.Config
}

type BackgroundProcessor interface {
	Begin(ProcessorConfig)
}
