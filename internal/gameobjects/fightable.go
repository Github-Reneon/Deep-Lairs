package gameobjects

type Fightable struct {
	Health         int
	MaxHealth      int
	BaseMaxHealth  int
	Attack         int
	BaseAttack     int
	Defense        int
	BaseDefense    int
	Mana           int
	MaxMana        int
	BaseMaxMana    int
	Stamina        int
	MaxStamina     int
	BaseMaxStamina int
	Speed          int
	BaseSpeed      int
	Int            int
	BaseInt        int
	InCombat       bool
	ItemStates     []*ItemState
	Image          string
}

type ItemState struct {
	Item     *Item
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
	f.ItemStates = []*ItemState{
		{
			Item:     CreateTestRingOfHealth(),
			Slot:     SLOT_RING,
			Equipped: false,
		},
		{
			Item:     CreateTestRingOfMana(),
			Slot:     SLOT_RING,
			Equipped: false,
		},
	}
}
