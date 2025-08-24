package gameobjects

type Fightable struct {
	Health         int
	MaxHealth      int
	baseMaxHealth  int
	Attack         int
	baseAttack     int
	Defense        int
	baseDefense    int
	Mana           int
	MaxMana        int
	baseMaxMana    int
	Stamina        int
	MaxStamina     int
	baseMaxStamina int
	XP             int
	maxXP          int
	Speed          int
	baseSpeed      int
	Int            int
	baseInt        int
	InCombat       bool
	Items          []*Item
	Equipped       []*Item
}

func (f *Fightable) InitFightable(health, attack, defense, mana, stamina, speed, intelligence int) {
	f.Health = health
	f.MaxHealth = health
	f.baseMaxHealth = health
	f.Attack = attack
	f.baseAttack = attack
	f.Defense = defense
	f.baseDefense = defense
	f.Mana = mana
	f.MaxMana = mana
	f.baseMaxMana = mana
	f.Stamina = stamina
	f.MaxStamina = stamina
	f.baseMaxStamina = stamina
	f.XP = 0
	f.maxXP = 100
	f.InCombat = false
	f.Speed = speed
	f.baseSpeed = speed
	f.Int = intelligence
	f.baseInt = intelligence
	f.Items = []*Item{
		CreateTestRingOfHealth(),
		CreateTestRingOfMana(),
	}
	f.Equipped = []*Item{}
}
