package view

import "sync"

const (
	QUIT    = "quit"
	PAUSE   = "pause"
	RUNNING = "Running"
)

type Event struct {
	Type  string
	mutex sync.Mutex
}

// Set the type of event that the user is going to do
func (event *Event) SetType(newType string) {
	event.mutex.Lock()
	defer event.mutex.Unlock()
	event.Type = newType
}

// Get the type of event the user is doing
func (event *Event) GetType() string {
	event.mutex.Lock()
	defer event.mutex.Unlock()
	return event.Type
}
