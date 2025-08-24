package main

import (
	"deep_lairs/internal/gameobjects"
	"deep_lairs/internal/protocol"
)

func InitPlace() *gameobjects.Place {
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
		JoiningMessage:   protocol.STUMBLES_IN,
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
	}
}
