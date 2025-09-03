package gameobjects

import (
	"deep_lairs/internal/protocol"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
)

var characterLocks sync.Map

type Character struct {
	ID             uuid.UUID
	Name           string
	LastName       string
	Location       *Place
	KnownLocations []*Place
	LocationId     string
	MessageQueue   []string
	Looked         bool
	Busy           bool
	changed        bool
	fighting       *IFightable
	CharacterFightable
}

type IFightable interface {
	BeAttacked(attackRoll int) int
	BeDamaged(damage int) int
}

type CharacterFightable struct {
	XP    int
	MaxXP int
	Level int
	Class string
	Fightable
}

func (c *Character) GetName() string {
	if c.Name == "" {
		// int to string
		return fmt.Sprintf("%d", c.ID)
	}
	return c.Name
}

func (c *Character) AddMessage(msg string) {
	muInterface, _ := characterLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	c.MessageQueue = append(c.MessageQueue, msg)
}

func (c *Character) ClearLastMessage() {
	muInterface, _ := characterLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	if len(c.MessageQueue) > 0 {
		c.MessageQueue = c.MessageQueue[1:]
	}
}

func (c *Character) ChangeLocation(newLocation *Place) {
	if !c.Busy {
		muInterface, _ := characterLocks.LoadOrStore(c.ID, &sync.Mutex{})
		mu := muInterface.(*sync.Mutex)
		mu.Lock()
		defer mu.Unlock()
		c.Location = newLocation
		c.Looked = false
	} else {
		c.AddMessage("You are busy and cannot change locations.")
	}
}

func (c *Character) AddKnownLocation(location *Place) {
	muInterface, _ := characterLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	if slices.Contains(c.KnownLocations, location) {
		return
	}
	c.KnownLocations = append(c.KnownLocations, location)
}

func (c *Character) IsKnownLocation(location *Place) bool {
	muInterface, _ := characterLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	for _, knownLocation := range c.KnownLocations {
		if knownLocation.ID == location.ID {
			log.Println("Found matching known location:", knownLocation.ID)
			return true
		}
	}
	return false
}

func (c *Character) GetState() string {
	muInterface, _ := characterLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	ret := struct {
		Type    string
		Name    string
		Health  string
		MP      string
		Stamina string
		XP      string
		Combat  bool
	}{
		Type:    protocol.STATE_TYPE_USER,
		Name:    c.GetName(),
		Health:  fmt.Sprintf("HP: %d/%d", c.Health, c.MaxHealth),
		MP:      fmt.Sprintf("MP: %d/%d", c.Mana, c.MaxMana),
		Stamina: fmt.Sprintf("ST: %d/%d", c.Stamina, c.MaxStamina),
		XP:      fmt.Sprintf("XP: %d/%d", c.XP, c.MaxXP),
		Combat:  c.InCombat,
	}
	jsonData, _ := json.Marshal(ret)
	return string(jsonData)
}

func (c *Character) SetBusyState(busy bool) {
	muInterface, _ := characterLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	c.Busy = busy
}

func (c *Character) Search() {
	if !c.Busy {
		c.SetBusyState(true)
		defer c.SetBusyState(false)
		searchFinishes := time.Now().Add(time.Duration(12-((rand.Intn(c.Speed)+1)+(rand.Intn(c.Int)+1))) * time.Second)
		for {
			c.AddMessage("Searching...")
			time.Sleep(1 * time.Second)
			if time.Now().After(searchFinishes) {
				break
			}
		}
		for direction, place := range c.Location.JoiningLocations {
			found := false
			if slices.Contains(c.KnownLocations, place) {
				c.AddMessage(fmt.Sprintf("%s is %s", place.Name, direction))
				found = true
			}
			if !found {
				if place.HiddenLocationMessage != "" {
					c.AddMessage(fmt.Sprintf(place.HiddenLocationMessage, direction))
				} else {
					c.AddMessage(fmt.Sprintf("%s is %s", place.Name, direction))
				}
			}
		}
	} else {
		c.AddMessage("You are busy.")
	}
}

func (c *Character) EquipItem(itemState *ItemState) {
	// lock
	muInterface, _ := characterLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	if itemState.Equipped {
		c.AddMessage("You are already wearing that item.")
		return
	}
	for _, ItemState := range c.ItemStates {
		if !ItemState.Equipped {
			continue
		}
		equippedItem := ItemState.Item
		if equippedItem.Slot == itemState.Slot {
			c.AddMessage(fmt.Sprintf("You are already wearing the following item: %s, which is in the same item slot.", equippedItem.Name))
			return
		}
	}
	itemState.Equipped = true
	c.changed = true
	c.AddMessage(fmt.Sprintf("You equip the %s.", itemState.Item.Name))
}

func (c *Character) UnequipItem(itemState *ItemState) {
	// lock
	muInterface, _ := characterLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	if itemState.Equipped {
		itemState.Equipped = false
		c.changed = true
		c.AddMessage(fmt.Sprintf("You unequip the %s.", itemState.Item.Name))
		return
	}
	c.AddMessage("You are not wearing that item.")
}

func (c *Character) StartCalcStatsHandler() {
	for {
		// Simulate stat calculation
		time.Sleep(time.Second)
		if c.changed {
			c.changed = false
			c.MaxHealth = c.BaseMaxHealth
			c.Mana = c.BaseMaxMana
			c.Stamina = c.BaseMaxStamina
			c.Attack = c.BaseAttack
			c.Defense = c.BaseDefense
			for _, ItemState := range c.ItemStates {
				if !ItemState.Equipped {
					continue
				}
				equippedItem := ItemState.Item
				switch equippedItem.BonusType {
				case BONUS_TYPE_ATTACK:
					c.Attack = c.BaseAttack + equippedItem.BonusAmount
				case BONUS_TYPE_DEFENSE:
					c.Defense = c.BaseDefense + equippedItem.BonusAmount
				case BONUS_TYPE_MANA:
					c.MaxMana = c.BaseMaxMana + equippedItem.BonusAmount
				case BONUS_TYPE_HEALTH:
					c.MaxHealth = c.BaseMaxHealth + equippedItem.BonusAmount
				}
			}
			// send back the state to the client
			c.AddMessage(c.GetState())
		}
	}
}

func (c *Character) Init(health, attack, defense, mana, stamina, speed, intelligence int) {
	// lock
	muInterface, _ := characterLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	c.CharacterFightable.InitFightable(health, attack, defense, mana, stamina, speed, intelligence)
	c.XP = 0
	c.MaxXP = 100
	c.Level = 1
	c.Image = "portrait_human_8.webp"
}

func (c *Character) Save() {
	// lock
	muInterface, _ := characterLocks.LoadOrStore(c.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	// save user state to json file in ./json/users/
	filePath := fmt.Sprintf("./json/users/%s.json", c.Name)
	data, err := json.Marshal(c)
	if err != nil {
		log.Println("Error marshalling user data:", err)
		return
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		log.Println("Error writing user data to file:", err)
	}
}
