package gameobjects

import "github.com/google/uuid"

const (
	BONUS_TYPE_HEALTH = iota
	BONUS_TYPE_ATTACK
	BONUS_TYPE_DEFENSE
	BONUS_TYPE_MANA
	BONUS_TYPE_STAMINA
	BONUS_TYPE_SPEED
)

const (
	SLOT_HEAD = iota
	SLOT_CHEST
	SLOT_LEGS
	SLOT_WEAPON
	SLOT_SHIELD
	SLOT_RING
	SLOT_AMULET
)

type Item struct {
	Id          string
	Name        string
	Description string
	Weight      int
	Value       int
	Type        string
	BonusType   int
	Slot        int
	BonusAmount int
	Tags        []string
}

func CreateTestRingOfHealth() *Item {
	return &Item{
		Id:          uuid.New().String(),
		Name:        "Ring of Health",
		Description: "A shiny ring that boosts your health.",
		Weight:      1,
		Value:       100,
		Type:        "ring",
		BonusType:   BONUS_TYPE_HEALTH,
		Slot:        SLOT_RING,
		BonusAmount: 5,
		Tags:        []string{"health", "ring"},
	}
}

func CreateTestRingOfMana() *Item {
	return &Item{
		Id:          uuid.New().String(),
		Name:        "Ring of Mana",
		Description: "A shiny ring that boosts your mana.",
		Weight:      1,
		Value:       100,
		Type:        "ring",
		BonusType:   BONUS_TYPE_MANA,
		Slot:        SLOT_RING,
		BonusAmount: 1,
		Tags:        []string{"mana", "ring"},
	}
}
