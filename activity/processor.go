package activity

import "sync"

type semaphore chan struct{}

type BackgroundProcessor interface {
	Begin(*sync.WaitGroup)
}
