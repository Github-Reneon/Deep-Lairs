package main

import (
	"deep_lairs/internal/gameobjects"
	"deep_lairs/internal/protocol"
)

func InitTavernPlace() *gameobjects.Place {
	return &gameobjects.Place{
		ID:          "tavern",
		Name:        "The Lofty Tavern",
		Description: "A cozy tavern with warm lighting.",
		TitleLook:   "You see a tavern wench. She is serving drinks.",
		Look: "There is a fireplace in the corner." +
			" Everyone seems to be enjoying themselves. Soft but lively music plays in the background.",
		LocationImage:    "tavern.webp",
		LookImage:        "drinks.webp",
		Users:            make(map[string]*gameobjects.User),
		Messages:         make([]string, 0),
		JoiningLocations: make(map[string]*gameobjects.Place),
		JoiningMessage:   protocol.JOINING_STUMBLES_IN,
		Jingles: []string{
			"Someone falls over.",
			"A bard starts playing a lively tune.",
			"A group of adventurers enters the tavern.",
			"Someone buys a round of drinks, and everyone cheers.",
			"The tavern wench sings a cheerful song.",
			"A mysterious figure in a hooded cloak sits in the corner, observing everyone... don't pay too much attention to them. They're looking for someone shorter than you.",
			"The fireplace crackles warmly, casting dancing shadows on the walls.",
			"Someone is getting handsy with the tavern wench... she hits back!",
			"A bar fight breaks out! Chairs are flying and people are shouting.",
			"A fledgling mage casts a spell, causing a small explosion in the corner!",
			"A cat jumps onto a table, knocking over drinks.",
			"You feel as if it's all going to be all right.",
		},
		LeavingMessage: protocol.LEAVING_STUMBLES_OUT,
	}
}

func InitTownSquarePlace() *gameobjects.Place {
	return &gameobjects.Place{
		ID:               "square",
		Name:             "The Town Square",
		Description:      "The bustling heart of the town, filled with merchants and townsfolk.",
		TitleLook:        "You see a vibrant square filled with people.",
		Look:             "Stalls line the streets, selling all manner of goods. A fountain bubbles in the center.",
		LocationImage:    "square.webp",
		LookImage:        "fountain.webp",
		Users:            make(map[string]*gameobjects.User),
		Messages:         make([]string, 0),
		JoiningLocations: make(map[string]*gameobjects.Place),
		JoiningMessage:   protocol.JOINING_MESSAGE,
		Jingles: []string{
			"Someone drops their coins, and a child quickly picks them up.",
			"A bard begins to play a lively tune, drawing a crowd.",
			"A merchant shouts about their wares, trying to attract customers.",
			"A group of adventurers discusses their latest quest.",
			"The fountain splashes water, creating a refreshing mist.",
		},
		LeavingMessage: protocol.LEAVING_MESSAGE,
	}
}
