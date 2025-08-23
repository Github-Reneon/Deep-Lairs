package gameobjects

import (
	"fmt"
	"log"
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

func (p *Place) RemoveMessage(msg string) {
	muInterface, _ := placeLocks.LoadOrStore(p.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	p.Messages = p.Messages[1:]
}

func (p *Place) RemoveUser(user *User) {
	muInterface, _ := placeLocks.LoadOrStore(p.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	delete(p.Users, user.ID)
}

func (p *Place) StartMessageHandler() {
	for {
		time.Sleep(time.Second * 1)
		if len(p.Messages) > 0 {
			log.Println("Sending messages to users in place:", p.Name)
			for _, message := range p.Messages {
				// Send the message to all users in the place
				for _, user := range p.Users {
					user.AddMessage(message)
				}
				p.RemoveMessage(message)
			}
		}
	}
}

func GetUser(world *World, id string) *User {
	for _, place := range world.Places {
		if user, ok := place.Users[id]; ok {
			return user
		}
	}
	return nil
}

func (w *World) StartJingleHandler() {
	for {
		time.Sleep(time.Minute)
		log.Println("Sending jingles to all places")
		for _, place := range w.Places {
			// random jingle add message
			jingle := place.Jingles[rand.Intn(len(place.Jingles))]
			place.AddMessage(jingle)
		}
	}
}
