package gameobjects

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var placeLocks sync.Map

type World struct {
	Places       map[string]*Place
	CurrentUsers int
}

func (w *World) AddShout(msg string) {
	for _, place := range w.Places {
		place.AddMessage(msg)
	}
}

type Place struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Users            map[string]*User
	QuickLook        string `json:"quick_look"`
	Look             string `json:"look"`
	LookImage        string `json:"look_image"`
	LocationImage    string `json:"location_image"`
	Messages         []string
	Jingles          []string `json:"jingles"`
	JoiningLocations map[string]*Place
}

func (p *Place) AddMessage(msg string) {
	muInterface, _ := placeLocks.LoadOrStore(p.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	p.Messages = append(p.Messages, msg)
}

func (p *Place) AddUser(user *User) {
	muInterface, _ := placeLocks.LoadOrStore(p.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	p.Users[user.ID] = user
}

func (p *Place) GetDirection(direction string) (*Place, error) {
	muInterface, _ := placeLocks.LoadOrStore(p.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	nextPlace, ok := p.JoiningLocations[direction]
	if !ok {
		return nil, fmt.Errorf("no place found in direction: %s", direction)
	}
	return nextPlace, nil
}

func GetUser(world *World, id string) *User {
	for _, place := range world.Places {
		if user, ok := place.Users[id]; ok {
			return user
		}
	}
	return nil
}

func (w *World) AddJingles() {
	for _, place := range w.Places {
		// random jingle add message
		jingle := place.Jingles[rand.Intn(len(place.Jingles))]
		place.AddMessage(jingle)
		time.Sleep(time.Second * time.Duration(rand.Intn(11)+5))
	}
}
