package gameobjects

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var placeLocks sync.Map

type Place struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Users            map[string]*User
	Look             string `json:"look"`
	TitleLook        string `json:"title_look"`
	LookImage        string `json:"look_image"`
	LocationImage    string `json:"location_image"`
	Messages         []string
	Jingles          []string `json:"jingles"`
	JoiningLocations map[string]*Place
	JoiningMessage   string  `json:"joining_message"`
	LeavingMessage   string  `json:"leaving_message"`
	Quests           []Quest `json:"quests"`
}

func (p *Place) AddUserMessage(msg string, user *User) {
	muInterface, _ := placeLocks.LoadOrStore(p.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	for _, u := range p.Users {
		if u != user {
			u.AddMessage(msg)
		}
	}
}

func (p *Place) AddMessage(msg string) {
	muInterface, _ := placeLocks.LoadOrStore(p.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	p.Messages = append(p.Messages, msg)
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

func (p *Place) RemoveUser(user *User, direction string) {
	defer p.AddUserMessage(fmt.Sprintf(p.LeavingMessage, user.GetName(), direction), user)
	muInterface, _ := placeLocks.LoadOrStore(p.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	delete(p.Users, user.ID)
}

func (p *Place) AddUser(user *User) {
	defer p.AddUserMessage(fmt.Sprintf(p.JoiningMessage, user.GetName()), user)
	muInterface, _ := placeLocks.LoadOrStore(p.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	p.Users[user.ID] = user
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

func (p *Place) StartCheckUsersHandler() {
	for {
		time.Sleep(time.Second * 5)
		if len(p.Users) > 0 {
			log.Println("Checking users in place:", p.Name)
			for _, user := range p.Users {
				if user.Location != p {
					log.Println("Removing user from place:", user.Name)
					p.RemoveUser(user, "poof")
					p.AddMessage(fmt.Sprintf("User %s has been removed from %s.", user.Name, p.Name))
				}
			}
		}
	}
}
