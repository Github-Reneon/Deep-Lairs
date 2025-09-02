package gameobjects

type Fightable struct {
	Health         int         `json:"-"`
	MaxHealth      int         `json:"health"`
	BaseMaxHealth  int         `json:"-"`
	Attack         int         `json:"-"`
	BaseAttack     int         `json:"attack"`
	Defense        int         `json:"-"`
	BaseDefense    int         `json:"defense"`
	Mana           int         `json:"-"`
	MaxMana        int         `json:"mana"`
	BaseMaxMana    int         `json:"-"`
	Stamina        int         `json:"-"`
	MaxStamina     int         `json:"stamina"`
	BaseMaxStamina int         `json:"-"`
	Speed          int         `json:"-"`
	BaseSpeed      int         `json:"speed"`
	Int            int         `json:"-"`
	BaseInt        int         `json:"intelligence"`
	InCombat       bool        `json:"-"`
	ItemStates     []ItemState `json:"item_states" gorm:"foreignKey:FightableID;constraint:onDelete:CASCADE,onUpdate:CASCADE"`
	Image          string      `json:"image"`
}

type ItemState struct {
	Item     *Item `gorm:"foreignKey:ItemID;constraint:onDelete:CASCADE,onUpdate:CASCADE"`
	ItemID   uint
	Slot     int
	Equipped bool
}

func (f *Fightable) InitFightable(health, attack, defense, mana, stamina, speed, intelligence int) {
	f.Health = health
	f.MaxHealth = health
	f.BaseMaxHealth = health
	f.Attack = attack
	f.BaseAttack = attack
	f.Defense = defense
	f.BaseDefense = defense
	f.Mana = mana
	f.MaxMana = mana
	f.BaseMaxMana = mana
	f.Stamina = stamina
	f.MaxStamina = stamina
	f.BaseMaxStamina = stamina
	f.InCombat = false
	f.Speed = speed
	f.BaseSpeed = speed
	f.Int = intelligence
	f.BaseInt = intelligence
	f.ItemStates = []ItemState{
		ItemState{
			Item:     CreateTestRingOfHealth(),
			Slot:     SLOT_RING,
			Equipped: false,
		},
		ItemState{
			Item:     CreateTestRingOfMana(),
			Slot:     SLOT_RING,
			Equipped: false,
		},
	}
}
