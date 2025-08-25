package gameobjects

type Enemy struct {
	Name     string `json:"name"`
	XPReward int    `json:"xp_reward"`
	Fightable
}
