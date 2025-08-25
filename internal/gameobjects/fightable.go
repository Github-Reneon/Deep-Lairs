package gameobjects

type Fightable struct {
	Health         int     `json:"health"`
	MaxHealth      int     `json:"max_health"`
	BaseMaxHealth  int     `json:"-"`
	Attack         int     `json:"attack"`
	BaseAttack     int     `json:"-"`
	Defense        int     `json:"defense"`
	BaseDefense    int     `json:"-"`
	Mana           int     `json:"mana"`
	MaxMana        int     `json:"max_mana"`
	BaseMaxMana    int     `json:"-"`
	Stamina        int     `json:"stamina"`
	MaxStamina     int     `json:"max_stamina"`
	BaseMaxStamina int     `json:"-"`
	XP             int     `json:"xp"`
	MaxXP          int     `json:"max_xp"`
	Speed          int     `json:"speed"`
	BaseSpeed      int     `json:"-"`
	Int            int     `json:"int"`
	BaseInt        int     `json:"-"`
	InCombat       bool    `json:"in_combat"`
	Items          []*Item `json:"items"`
	Equipped       []*Item `json:"equipped"`
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
	f.XP = 0
	f.MaxXP = 100
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
