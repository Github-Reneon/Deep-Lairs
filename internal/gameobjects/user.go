package gameobjects

import (
	"log"
	"slices"
	"sync"
)

var userLocks sync.Map

type User struct {
	ID             string
	Name           string
	Location       *Place
	MessageQueue   []string
	Looked         bool
	KnownLocations []*Place
}

func (u *User) GetName() string {
	if u.Name == "" {
		return u.ID
	}
	return u.Name
}

func (u *User) AddMessage(msg string) {
	muInterface, _ := userLocks.LoadOrStore(u.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	u.MessageQueue = append(u.MessageQueue, msg)
}

func (u *User) ClearLastMessage() {
	muInterface, _ := userLocks.LoadOrStore(u.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	if len(u.MessageQueue) > 0 {
		u.MessageQueue = u.MessageQueue[1:]
	}
}

func (u *User) ChangeLocation(newLocation *Place) {
	muInterface, _ := userLocks.LoadOrStore(u.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	u.Location = newLocation
	u.Looked = false
}

func (u *User) AddKnownLocation(location *Place) {
	muInterface, _ := userLocks.LoadOrStore(u.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	if slices.Contains(u.KnownLocations, location) {
		return
	}
	u.KnownLocations = append(u.KnownLocations, location)
}

func (u *User) IsKnownLocation(location *Place) bool {
	muInterface, _ := userLocks.LoadOrStore(u.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	for _, knownLocation := range u.KnownLocations {
		if knownLocation.ID == location.ID {
			log.Println("Found matching known location:", knownLocation.ID)
			return true
		}
	}
	return false
}
