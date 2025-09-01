package gameobjects

import (
	"log"
	"math/rand"
	"time"
)

type World struct {
	Places            map[string]*Place
	CurrentCharacters int
}

func (w *World) AddShout(msg string) {
	for _, place := range w.Places {
		place.AddMessage(msg)
	}
}

func GetCharacter(world *World, id string) *Character {
	for _, place := range world.Places {
		if user, ok := place.Characters[id]; ok {
			return user
		}
	}
	return nil
}

func (w *World) StartJingleHandler() {
	for {
		time.Sleep(time.Minute)
		log.Println("Sending jingles to all places")
		for _, place := range w.Places {
			// random jingle add message
			jingle := place.Jingles[rand.Intn(len(place.Jingles))]
			place.AddMessage(jingle)
		}
	}
}
