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
	Items          []*Item     `json:"-"`
	ItemStates     []ItemState `json:"items"`
	Equipped       []*Item     `json:"-"`
}

type ItemState struct {
	ItemId   string `json:"item_id"`
	Equipped bool   `json:"equipped"`
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
	f.Items = []*Item{
		CreateTestRingOfHealth(),
		CreateTestRingOfMana(),
	}
	f.Equipped = []*Item{}
}
