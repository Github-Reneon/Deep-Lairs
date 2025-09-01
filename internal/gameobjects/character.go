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
)

var userLocks sync.Map

type Character struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	LastName         string      `json:"last_name"`
	Location         *Place      `json:"-"`
	LocationId       string      `json:"location"`
	MessageQueue     []string    `json:"-"`
	Looked           bool        `json:"-"`
	KnownLocations   []*Place    `json:"-"`
	KnownLocationIds []string    `json:"known_locations"`
	Busy             bool        `json:"-"`
	changed          bool        `json:"-"`
	fighting         *IFightable `json:"-"`
	UserFightable
}

type IFightable interface {
	BeAttacked(attackRoll int) int
	BeDamaged(damage int) int
}

type UserFightable struct {
	XP    int `json:"xp"`
	MaxXP int `json:"max_xp"`
	Level int `json:"level"`
	Fightable
}

func (c *Character) GetName() string {
	if c.Name == "" {
		return c.ID
	}
	return c.Name
}

func (c *Character) AddMessage(msg string) {
	muInterface, _ := userLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	c.MessageQueue = append(c.MessageQueue, msg)
}

func (c *Character) ClearLastMessage() {
	muInterface, _ := userLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	if len(c.MessageQueue) > 0 {
		c.MessageQueue = c.MessageQueue[1:]
	}
}

func (c *Character) ChangeLocation(newLocation *Place) {
	if !c.Busy {
		muInterface, _ := userLocks.LoadOrStore(c.ID, &sync.Mutex{})
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
	muInterface, _ := userLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	if slices.Contains(c.KnownLocations, location) {
		return
	}
	c.KnownLocations = append(c.KnownLocations, location)
}

func (c *Character) IsKnownLocation(location *Place) bool {
	muInterface, _ := userLocks.LoadOrStore(c.ID, &sync.Mutex{})
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
	muInterface, _ := userLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	ret := struct {
		Type    string `json:"type"`
		Name    string `json:"name"`
		Health  string `json:"hp"`
		MP      string `json:"mp"`
		Stamina string `json:"stamina"`
		XP      string `json:"xp"`
		Combat  bool   `json:"combat"`
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
	muInterface, _ := userLocks.LoadOrStore(c.ID, &sync.Mutex{})
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

func (c *Character) EquipItem(item *Item) {
	// lock
	defer c.AddMessage(fmt.Sprintf("You equip the %s.", item.Name))
	muInterface, _ := userLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	c.changed = true
	c.Equipped = append(c.Equipped, item)
}

func (c *Character) UnequipItem(item *Item) {
	// lock
	defer c.AddMessage(fmt.Sprintf("You unequip the %s.", item.Name))
	muInterface, _ := userLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	for i, equippedItem := range c.Equipped {
		if equippedItem == item {
			c.Equipped = append(c.Equipped[:i], c.Equipped[i+1:]...)
			c.changed = true
			return
		}
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
			for _, equippedItem := range c.Equipped {
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
	defer c.SetIds()
	muInterface, _ := userLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	c.UserFightable.InitFightable(health, attack, defense, mana, stamina, speed, intelligence)
	c.XP = 0
	c.MaxXP = 100
	c.Level = 1
	c.Image = "portrait_human_8.webp"
}

// SetIds assigns unique IDs to the user and their items.
func (c *Character) SetIds() {
	muInterface, _ := userLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	c.KnownLocationIds = make([]string, len(c.KnownLocations))
	for i, location := range c.KnownLocations {
		c.KnownLocationIds[i] = location.ID
	}
	c.ItemStates = make([]ItemState, len(c.Items))
	for i, item := range c.Items {
		c.ItemStates[i] = ItemState{
			ItemId:   item.Id,
			Equipped: slices.Contains(c.Equipped, item),
		}
	}
}

func (c *Character) StartSetIdsHandler() {
	// Start a goroutine to set IDs
	for {
		time.Sleep(time.Second)
		c.SetIds()
	}
}

func (c *Character) Save() {
	// lock
	muInterface, _ := userLocks.LoadOrStore(c.Name, &sync.Mutex{})
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
