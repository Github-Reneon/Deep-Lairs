package gameobjects

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var placeLocks sync.Map

type Place struct {
	ID                    string                `json:"id"`
	Name                  string                `json:"name"`
	Description           string                `json:"description"`
	Characters            map[string]*Character `json:"-"`
	Look                  string                `json:"look"`
	TitleLook             string                `json:"title_look"`
	LookImage             string                `json:"look_image"`
	LocationImage         string                `json:"location_image"`
	Messages              []string              `json:"-"`
	Jingles               []string              `json:"jingles"`
	JoiningLocations      map[string]*Place     `json:"-"`
	JoiningLocationIds    map[string]string     `json:"joining_location_ids"`
	JoiningMessage        string                `json:"joining_message"`
	LeavingMessage        string                `json:"leaving_message"`
	Quests                []Quest               `json:"quests,omitempty"`
	HiddenLocationMessage string                `json:"hidden_location_message,omitempty"`
}

func (p *Place) AddCharacterMessage(msg string, character *Character) {
	muInterface, _ := placeLocks.LoadOrStore(p.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	for _, u := range p.Characters {
		if u != character {
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

func (p *Place) RemoveCharacter(character *Character, direction string) {
	defer p.AddCharacterMessage(fmt.Sprintf(p.LeavingMessage, character.GetName(), direction), character)
	muInterface, _ := placeLocks.LoadOrStore(p.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	delete(p.Characters, character.ID)
}

func (p *Place) AddCharacter(character *Character) {
	defer p.AddCharacterMessage(fmt.Sprintf(p.JoiningMessage, character.GetName()), character)
	muInterface, _ := placeLocks.LoadOrStore(p.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	p.Characters[character.ID] = character
}

func (p *Place) StartMessageHandler() {
	for {
		time.Sleep(time.Second * 1)
		if len(p.Messages) > 0 {
			log.Println("Sending messages to users in place:", p.Name)
			for _, message := range p.Messages {
				// Send the message to all users in the place
				for _, user := range p.Characters {
					user.AddMessage(message)
				}
				p.RemoveMessage(message)
			}
		}
	}
}

func (p *Place) StartCheckCharactersHandler() {
	for {
		time.Sleep(time.Second * 5)
		if len(p.Characters) > 0 {
			for _, user := range p.Characters {
				if user.Location != p {
					p.RemoveCharacter(user, "poof")
					p.AddMessage(fmt.Sprintf("User %s has been removed from %s.", user.Name, p.Name))
				}
			}
		}
	}
}

func (p *Place) Init(world *World) {
	p.JoiningLocations = make(map[string]*Place)
	p.Characters = make(map[string]*Character)
	for dir, loc := range p.JoiningLocationIds {
		log.Println("Initializing joining location:", dir, loc)
		log.Println(world.Places[loc].ID, world.Places[loc].Name)
		p.JoiningLocations[dir] = world.Places[loc]
	}
}
