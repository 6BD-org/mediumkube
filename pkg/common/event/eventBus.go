package event

import (
	"mediumkube/pkg/models"
	"sync"
)

var (
	bus *eventBus = nil
	mux           = sync.Mutex{}
)

type DomainEvent struct {
	d         *models.Domain
	timestemp int64
}

type eventBus struct {
	DomainUpdate chan DomainEvent
}

func GetEventBus() *eventBus {
	if bus != nil {
		return bus
	}
	mux.Lock()
	defer mux.Unlock()
	if bus == nil {
		bus = &eventBus{
			DomainUpdate: make(chan DomainEvent, 32),
		}
	}
	return bus
}
