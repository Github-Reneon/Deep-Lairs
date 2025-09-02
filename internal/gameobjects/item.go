package gameobjects

import (
	"gorm.io/gorm"
)

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
	gorm.Model
	Name        string
	Description string
	Weight      int
	Value       int
	Type        string
	BonusType   int
	Slot        int
	BonusAmount int
	Tags        []*Tag `gorm:"many2many:item_tags;ondelete:CASCADE;onupdate:CASCADE"`
}

type Tag struct {
	gorm.Model
	Name string
}

func CreateTestRingOfHealth() *Item {
	return &Item{
		Name:        "Ring of Health",
		Description: "A shiny ring that boosts your health.",
		Weight:      1,
		Value:       100,
		Type:        "ring",
		BonusType:   BONUS_TYPE_HEALTH,
		Slot:        SLOT_RING,
		BonusAmount: 5,
		Tags: []*Tag{
			&Tag{
				Name: "health",
			},
			&Tag{
				Name: "ring",
			},
		},
	}
}

func CreateTestRingOfMana() *Item {
	return &Item{
		Name:        "Ring of Mana",
		Description: "A shiny ring that boosts your mana.",
		Weight:      1,
		Value:       100,
		Type:        "ring",
		BonusType:   BONUS_TYPE_MANA,
		Slot:        SLOT_RING,
		BonusAmount: 1,
		Tags: []*Tag{
			&Tag{
				Name: "mana",
			},
			&Tag{
				Name: "ring",
			},
		},
	}
}
