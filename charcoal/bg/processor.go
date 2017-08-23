package bg

import "sync"

type semaphore chan struct{}

// BackgroundProcessor defines an interface for continuous processes that run in the background.
type BackgroundProcessor interface {
	Begin(*sync.WaitGroup)
}
