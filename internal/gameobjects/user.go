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

type User struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
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
	if !u.Busy {
		muInterface, _ := userLocks.LoadOrStore(u.ID, &sync.Mutex{})
		mu := muInterface.(*sync.Mutex)
		mu.Lock()
		defer mu.Unlock()
		u.Location = newLocation
		u.Looked = false
	} else {
		u.AddMessage("You are busy and cannot change locations.")
	}
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

func (u *User) GetState() string {
	muInterface, _ := userLocks.LoadOrStore(u.ID, &sync.Mutex{})
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
		Name:    u.GetName(),
		Health:  fmt.Sprintf("HP: %d/%d", u.Health, u.MaxHealth),
		MP:      fmt.Sprintf("MP: %d/%d", u.Mana, u.MaxMana),
		Stamina: fmt.Sprintf("ST: %d/%d", u.Stamina, u.MaxStamina),
		XP:      fmt.Sprintf("XP: %d/%d", u.XP, u.MaxXP),
		Combat:  u.InCombat,
	}
	jsonData, _ := json.Marshal(ret)
	return string(jsonData)
}

func (u *User) SetBusyState(busy bool) {
	muInterface, _ := userLocks.LoadOrStore(u.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	u.Busy = busy
}

func (u *User) Search() {
	if !u.Busy {
		u.SetBusyState(true)
		defer u.SetBusyState(false)
		searchFinishes := time.Now().Add(time.Duration(12-((rand.Intn(u.Speed)+1)+(rand.Intn(u.Int)+1))) * time.Second)
		for {
			u.AddMessage("Searching...")
			time.Sleep(1 * time.Second)
			if time.Now().After(searchFinishes) {
				break
			}
		}
		for direction, place := range u.Location.JoiningLocations {
			found := false
			if slices.Contains(u.KnownLocations, place) {
				u.AddMessage(fmt.Sprintf("%s is %s", place.Name, direction))
				found = true
			}
			if !found {
				if place.HiddenLocationMessage != "" {
					u.AddMessage(fmt.Sprintf(place.HiddenLocationMessage, direction))
				} else {
					u.AddMessage(fmt.Sprintf("%s is %s", place.Name, direction))
				}
			}
		}
	} else {
		u.AddMessage("You are busy.")
	}
}

func (u *User) EquipItem(item *Item) {
	// lock
	defer u.AddMessage(fmt.Sprintf("You equip the %s.", item.Name))
	muInterface, _ := userLocks.LoadOrStore(u.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	u.changed = true
	u.Equipped = append(u.Equipped, item)
}

func (u *User) UnequipItem(item *Item) {
	// lock
	defer u.AddMessage(fmt.Sprintf("You unequip the %s.", item.Name))
	muInterface, _ := userLocks.LoadOrStore(u.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	for i, equippedItem := range u.Equipped {
		if equippedItem == item {
			u.Equipped = append(u.Equipped[:i], u.Equipped[i+1:]...)
			u.changed = true
			return
		}
	}
	u.AddMessage("You are not wearing that item.")
}

func (u *User) StartCalcStatsHandler() {
	for {
		// Simulate stat calculation
		time.Sleep(time.Second)
		if u.changed {
			u.changed = false
			u.MaxHealth = u.BaseMaxHealth
			u.Mana = u.BaseMaxMana
			u.Stamina = u.BaseMaxStamina
			u.Attack = u.BaseAttack
			u.Defense = u.BaseDefense
			for _, equippedItem := range u.Equipped {
				switch equippedItem.BonusType {
				case BONUS_TYPE_ATTACK:
					u.Attack = u.BaseAttack + equippedItem.BonusAmount
				case BONUS_TYPE_DEFENSE:
					u.Defense = u.BaseDefense + equippedItem.BonusAmount
				case BONUS_TYPE_MANA:
					u.MaxMana = u.BaseMaxMana + equippedItem.BonusAmount
				case BONUS_TYPE_HEALTH:
					u.MaxHealth = u.BaseMaxHealth + equippedItem.BonusAmount
				}
			}
			// send back the state to the client
			u.AddMessage(u.GetState())
		}
	}
}

func (u *User) Init(health, attack, defense, mana, stamina, speed, intelligence int) {
	// lock
	defer u.SetIds()
	muInterface, _ := userLocks.LoadOrStore(u.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	u.UserFightable.InitFightable(health, attack, defense, mana, stamina, speed, intelligence)
	u.XP = 0
	u.MaxXP = 100
	u.Level = 1
	u.Image = "portrait_human_8.webp"
}

// SetIds assigns unique IDs to the user and their items.
func (u *User) SetIds() {
	muInterface, _ := userLocks.LoadOrStore(u.ID, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	u.KnownLocationIds = make([]string, len(u.KnownLocations))
	for i, location := range u.KnownLocations {
		u.KnownLocationIds[i] = location.ID
	}
	u.ItemStates = make([]ItemState, len(u.Items))
	for i, item := range u.Items {
		u.ItemStates[i] = ItemState{
			ItemId:   item.Id,
			Equipped: slices.Contains(u.Equipped, item),
		}
	}
}

func (u *User) StartSetIdsHandler() {
	// Start a goroutine to set IDs
	for {
		time.Sleep(time.Second)
		u.SetIds()
	}
}

func (u *User) Save() {
	// lock
	muInterface, _ := userLocks.LoadOrStore(u.Name, &sync.Mutex{})
	mu := muInterface.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	// save user state to json file in ./json/users/
	filePath := fmt.Sprintf("./json/users/%s.json", u.Name)
	data, err := json.Marshal(u)
	if err != nil {
		log.Println("Error marshalling user data:", err)
		return
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		log.Println("Error writing user data to file:", err)
	}
}
