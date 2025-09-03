package main

import (
	"deep_lairs/internal/gameobjects"
	"deep_lairs/internal/protocol"
	"fmt"
	"slices"
	"strings"
)

func CharacterLaugh(character *gameobjects.Character) {
	character.AddMessage(fmt.Sprintf(protocol.LOL, character.GetName()))
	for _, u := range character.Location.Characters {
		if u != character {
			u.AddMessage(fmt.Sprintf(protocol.LOL, character.GetName()))
		}
	}
}

func CharacterShout(splitMsg []string, character *gameobjects.Character) {
	if len(splitMsg) < 2 {
		character.AddMessage("Usage: shout <message>")
	} else {
		character.AddMessage(fmt.Sprintf(protocol.SHOUT, character.GetName(), strings.ToUpper(strings.Join(splitMsg[1:], " "))))
		// replace later with adding the message to each location
		for _, u := range character.Location.Characters {
			if u != character {
				u.AddMessage(fmt.Sprintf(protocol.SHOUT, character.GetName(), strings.ToUpper(strings.Join(splitMsg[1:], " "))))
			}
		}
	}
}

func CharacterSay(splitMsg []string, character *gameobjects.Character) {
	if len(splitMsg) < 2 {
		character.AddMessage("Usage: say <message>")
	} else {
		message := fmt.Sprintf(
			protocol.SAY,
			"/img/portraits/"+character.Image,
			character.GetName(),
			strings.ToUpper(splitMsg[1][:1])+strings.Join(splitMsg[1:], " ")[1:],
		)
		character.AddMessage(message)
		for _, u := range character.Location.Characters {
			if u != character {
				u.AddMessage(message)
			}
		}
	}
}

func CharacterLook(splitMsg []string, character *gameobjects.Character) {
	if len(splitMsg) > 1 {
		character.AddMessage(fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, "look "+strings.Join(splitMsg[1:], " ")))
	} else {
		b := strings.Builder{}
		b.WriteString(fmt.Sprintf(protocol.LOOK, character.Location.TitleLook, character.Location.Look, character.Location.LookImage))
		if len(character.Location.JoiningLocations) > 0 {
			for direction, place := range character.Location.JoiningLocations {
				if slices.Contains(character.KnownLocations, place) && place.HiddenLocationMessage != "" {
					b.WriteString(fmt.Sprintf(place.HiddenLocationMessage, direction))
				}
			}
		}
		character.AddMessage(b.String())
		lookCharacters(character)
	}
}

func lookCharacters(character *gameobjects.Character) {
	characters := []string{}
	for _, foundCharacter := range character.Location.Characters {
		if foundCharacter != character {
			characters = append(characters, foundCharacter.GetName())
		}
	}
	if len(characters) >= 1 {
		if len(characters) >= 10 {
			character.AddMessage(fmt.Sprintf("You see many adventurers here. %d in total.", len(characters)))
		} else {
			b := strings.Builder{}
			b.WriteString("You see the following adventurers here:<br class=\"my-2\">")
			b.WriteString("<span class=\"inline-grid grid-cols-5 gap-4\">")
			for _, u := range characters {
				b.WriteString(fmt.Sprintf("<span class=\"p-2\">%s</span>", u))
			}
			b.WriteString("</span>")
			character.AddMessage(b.String())
		}
	} else {
		character.AddMessage("You don't see any other adventurers here.")
	}
}

func CharacterQuickLook(splitMsg []string, character *gameobjects.Character) {
	if len(splitMsg) > 1 {
		character.AddMessage(fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, "quick_look "+strings.Join(splitMsg[1:], " ")))
	} else {
		b := strings.Builder{}
		b.WriteString(fmt.Sprintf(protocol.LOOK_NO_IMAGE, character.Location.TitleLook, character.Location.Look))
		b.WriteString(" ")
		if len(character.Location.JoiningLocations) > 0 {
			for direction, place := range character.Location.JoiningLocations {
				if slices.Contains(character.KnownLocations, place) && place.HiddenLocationMessage != "" {
					b.WriteString(fmt.Sprintf(place.HiddenLocationMessage, direction))
				}
			}
		}
		character.AddMessage(b.String())
	}
	lookCharacters(character)
}

func CharacterGo(splitMsg []string, character *gameobjects.Character) (bool, error) {
	if len(splitMsg) < 2 {
		character.AddMessage("Usage: go <direction>")
		return false, fmt.Errorf("no direction provided")
	}
	direction := strings.ToLower(splitMsg[1])
	switch direction {
	case "n":
		direction = "north"
	case "s":
		direction = "south"
	case "e":
		direction = "east"
	case "w":
		direction = "west"
	case "d":
		direction = "down"
	case "u":
		direction = "up"
	case "i":
		direction = "in"
	case "o":
		direction = "out"
	}
	if newLocation, ok := character.Location.JoiningLocations[direction]; ok {
		knownLocation := character.IsKnownLocation(newLocation)
		character.Location.RemoveCharacter(character, direction)
		newLocation.AddCharacter(character)
		character.ChangeLocation(newLocation)
		character.AddMessage(fmt.Sprintf("You go %s.", direction))
		character.AddKnownLocation(newLocation)
		return knownLocation, nil
	} else {
		character.AddMessage(fmt.Sprintf("You can't go %s.", direction))
		return false, fmt.Errorf("no location in direction: %s", direction)
	}
}

func CharacterWhere(splitMsg []string, character *gameobjects.Character) {
	character.AddMessage(fmt.Sprintf("You are in %s<br>%s", character.Location.Name, character.Location.Description))
}

func CharacterSearch(splitMsg []string, character *gameobjects.Character) {
	if len(splitMsg) > 1 {
		character.AddMessage(fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, "search "+strings.Join(splitMsg[1:], " ")))
	} else {
		go character.Search()
	}
}

func CharacterJoin(character *gameobjects.Character) {
	CharacterWhere([]string{"w"}, character)
	character.AddMessage(fmt.Sprintf(protocol.IMAGE, character.Location.LocationImage))
	CharacterQuickLook([]string{"l"}, character)
}

func CharacterQuestBoard(splitMsg []string, character *gameobjects.Character) {
	if len(splitMsg) > 1 {
		character.AddMessage(fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, "questboard "+strings.Join(splitMsg[1:], " ")))
	} else {
		if len(character.Location.Quests) == 0 {
			character.AddMessage("There is no quest board here.")
		} else {
			b := strings.Builder{}
			b.WriteString("You see a quest board with the following quests:<br><span class=\"inline-grid grid-cols-5 gap-4\">")
			// List available quests
			for a, quest := range character.Location.Quests {
				b.WriteString(fmt.Sprintf("<span class=\"p-2 bg-gray-900 rounded-lg\">%s (%d)</span>", quest.Name, a+1))
			}
			b.WriteString("</span>")
			character.AddMessage(b.String())
		}
	}
}

// TODO
func SendCharacterState(character *gameobjects.Character) {
	character.AddMessage(character.GetState())
}

func CharacterInventory(character *gameobjects.Character) {
	if len(character.ItemStates) == 0 {
		character.AddMessage("Your inventory is empty.")
	} else {
		b := strings.Builder{}
		b.WriteString("You have the following items in your inventory:<br><span class=\"inline-grid grid-cols-5 gap-4\">")
		for _, ItemState := range character.ItemStates {
			if ItemState.Equipped {
				b.WriteString(fmt.Sprintf("<span class=\"p-2 bg-green-700 rounded-lg\">%s (equipped)</span>", ItemState.Item.Name))
			} else {
				b.WriteString(fmt.Sprintf("<span class=\"p-2 bg-gray-900 rounded-lg\">%s</span>", ItemState.Item.Name))
			}
		}
		b.WriteString("</span>")
		character.AddMessage(b.String())
	}
}

func CharacterEquip(splitMsg []string, character *gameobjects.Character) {
	if len(splitMsg) < 2 {
		character.AddMessage("Usage: equip <item>")
		return
	}
	itemName := strings.Join(splitMsg[1:], " ")
	itemState, notFound := findItemState(character, itemName)
	if notFound {
		character.AddMessage("Item not found.")
		return
	}
	if itemState.Equipped {
		character.AddMessage("You are already wearing that item.")
		return
	}
	for _, characterItemState := range character.ItemStates {
		if !characterItemState.Equipped {
			continue
		}
		equippedItem := characterItemState.Item
		if equippedItem.Slot == itemState.Slot {
			character.AddMessage(fmt.Sprintf("You are already wearing the following item: %s, which is in the same item slot.", equippedItem.Name))
			return
		}
	}
	character.EquipItem(itemState)
}

func findItemState(character *gameobjects.Character, itemName string) (*gameobjects.ItemState, bool) {
	var foundItemState *gameobjects.ItemState
	for _, ItemState := range character.ItemStates {
		if ItemState.Item.Name == itemName {
			foundItemState = ItemState
			break
		}
	}
	if foundItemState == nil {
		// search tags
		for _, ItemState := range character.ItemStates {
			for _, tag := range ItemState.Item.Tags {
				if tag.Name == itemName {
					foundItemState = ItemState
					break
				}
			}
		}
	}
	if foundItemState == nil {
		return nil, true
	}
	return foundItemState, false
}

func CharacterUnequip(splitMsg []string, character *gameobjects.Character) {
	if len(splitMsg) < 2 {
		character.AddMessage("Usage: unequip <item>")
		return
	}
	itemName := strings.Join(splitMsg[1:], " ")
	itemState, shouldReturn := findItemState(character, itemName)
	if shouldReturn {
		return
	}
	if !itemState.Equipped {
		character.AddMessage("You are not wearing that item.")
		return
	}
	character.UnequipItem(itemState)
}

func CharacterDo(splitMsg []string, character *gameobjects.Character) {
	if len(splitMsg) < 2 {
		character.AddMessage("Usage: do <action>")
		return
	}
	actionName := strings.Join(splitMsg[1:], " ")
	character.Location.AddMessage(fmt.Sprintf(protocol.DO, character.GetName(), actionName))
}
