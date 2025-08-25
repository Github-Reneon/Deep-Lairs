package main

import (
	"deep_lairs/internal/gameobjects"
	"deep_lairs/internal/protocol"
	"fmt"
	"slices"
	"strings"
)

func UserLaugh(user *gameobjects.User) {
	user.AddMessage(fmt.Sprintf(protocol.LOL, user.GetName()))
	for _, u := range user.Location.Users {
		if u != user {
			u.AddMessage(fmt.Sprintf(protocol.LOL, user.GetName()))
		}
	}
}

func UserShout(splitMsg []string, user *gameobjects.User) {
	if len(splitMsg) < 2 {
		user.AddMessage("Usage: shout <message>")
	} else {
		user.AddMessage(fmt.Sprintf(protocol.SHOUT, user.GetName(), strings.ToUpper(strings.Join(splitMsg[1:], " "))))
		// replace later with adding the message to each location
		for _, u := range user.Location.Users {
			if u != user {
				u.AddMessage(fmt.Sprintf(protocol.SHOUT, user.GetName(), strings.ToUpper(strings.Join(splitMsg[1:], " "))))
			}
		}
	}
}

func UserSay(splitMsg []string, user *gameobjects.User) {
	if len(splitMsg) < 2 {
		user.AddMessage("Usage: say <message>")
	} else {
		message := fmt.Sprintf(
			protocol.SAY,
			user.GetName(),
			strings.ToUpper(splitMsg[1][:1])+strings.Join(splitMsg[1:], " ")[1:],
		)
		user.AddMessage(message)
		for _, u := range user.Location.Users {
			if u != user {
				u.AddMessage(message)
			}
		}
	}
}

func UserLook(splitMsg []string, user *gameobjects.User) {
	if len(splitMsg) > 1 {
		user.AddMessage(fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, "look "+strings.Join(splitMsg[1:], " ")))
	} else {
		b := strings.Builder{}
		b.WriteString(fmt.Sprintf(protocol.LOOK, user.Location.TitleLook, user.Location.Look, user.Location.LookImage))
		if len(user.Location.JoiningLocations) > 0 {
			for direction, place := range user.Location.JoiningLocations {
				if slices.Contains(user.KnownLocations, place) && place.HiddenLocationMessage != "" {
					b.WriteString(fmt.Sprintf(place.HiddenLocationMessage, direction))
				}
			}
		}
		user.AddMessage(b.String())
		lookUsers(user)
	}
}

func lookUsers(user *gameobjects.User) {
	users := []string{}
	for _, foundUser := range user.Location.Users {
		if foundUser != user {
			users = append(users, foundUser.GetName())
		}
	}
	if len(users) >= 1 {
		if len(users) >= 10 {
			user.AddMessage(fmt.Sprintf("You see many adventurers here. %d in total.", len(users)))
		} else {
			b := strings.Builder{}
			b.WriteString("You see the following adventurers here:<br class=\"my-2\">")
			b.WriteString("<span class=\"inline-grid grid-cols-5 gap-4\">")
			for _, u := range users {
				b.WriteString(fmt.Sprintf("<span class=\"p-2\">%s</span>", u))
			}
			b.WriteString("</span>")
			user.AddMessage(b.String())
		}
	} else {
		user.AddMessage("You don't see any other adventurers here.")
	}
}

func UserQuickLook(splitMsg []string, user *gameobjects.User) {
	if len(splitMsg) > 1 {
		user.AddMessage(fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, "quick_look "+strings.Join(splitMsg[1:], " ")))
	} else {
		b := strings.Builder{}
		b.WriteString(fmt.Sprintf(protocol.LOOK_NO_IMAGE, user.Location.TitleLook, user.Location.Look))
		b.WriteString(" ")
		if len(user.Location.JoiningLocations) > 0 {
			for direction, place := range user.Location.JoiningLocations {
				if slices.Contains(user.KnownLocations, place) && place.HiddenLocationMessage != "" {
					b.WriteString(fmt.Sprintf(place.HiddenLocationMessage, direction))
				}
			}
		}
		user.AddMessage(b.String())
	}
	lookUsers(user)
}

func UserGo(splitMsg []string, user *gameobjects.User) (bool, error) {
	if len(splitMsg) < 2 {
		user.AddMessage("Usage: go <direction>")
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
	if newLocation, ok := user.Location.JoiningLocations[direction]; ok {
		knownLocation := user.IsKnownLocation(newLocation)
		user.Location.RemoveUser(user, direction)
		newLocation.AddUser(user)
		user.ChangeLocation(newLocation)
		user.AddMessage(fmt.Sprintf("You go %s.", direction))
		user.AddKnownLocation(newLocation)
		return knownLocation, nil
	} else {
		user.AddMessage(fmt.Sprintf("You can't go %s.", direction))
		return false, fmt.Errorf("no location in direction: %s", direction)
	}
}

func UserWhere(splitMsg []string, user *gameobjects.User) {
	user.AddMessage(fmt.Sprintf("You are in %s<br>%s", user.Location.Name, user.Location.Description))
}

func UserSearch(splitMsg []string, user *gameobjects.User) {
	if len(splitMsg) > 1 {
		user.AddMessage(fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, "search "+strings.Join(splitMsg[1:], " ")))
	} else {
		go user.Search()
	}
}

func UserJoin(user *gameobjects.User) {
	UserWhere([]string{"w"}, user)
	user.AddMessage(fmt.Sprintf(protocol.IMAGE, user.Location.LocationImage))
	UserQuickLook([]string{"l"}, user)
}

func UserQuestBoard(splitMsg []string, user *gameobjects.User) {
	if len(splitMsg) > 1 {
		user.AddMessage(fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, "questboard "+strings.Join(splitMsg[1:], " ")))
	} else {
		if len(user.Location.Quests) == 0 {
			user.AddMessage("There is no quest board here.")
		} else {
			b := strings.Builder{}
			b.WriteString("You see a quest board with the following quests:<br><span class=\"inline-grid grid-cols-5 gap-4\">")
			// List available quests
			for a, quest := range user.Location.Quests {
				b.WriteString(fmt.Sprintf("<span class=\"p-2 bg-gray-900 rounded-lg\">%s (%d)</span>", quest.Name, a+1))
			}
			b.WriteString("</span>")
			user.AddMessage(b.String())
		}
	}
}

// TODO
func SendUserState(user *gameobjects.User) {
	user.AddMessage(user.GetState())
}

func UserInventory(user *gameobjects.User) {
	if len(user.Items) == 0 {
		user.AddMessage("Your inventory is empty.")
	} else {
		b := strings.Builder{}
		b.WriteString("You have the following items in your inventory:<br><span class=\"inline-grid grid-cols-5 gap-4\">")
		for _, item := range user.Items {
			if slices.Contains(user.Equipped, item) {
				b.WriteString(fmt.Sprintf("<span class=\"p-2 bg-green-700 rounded-lg\">%s (equipped)</span>", item.Name))
			} else {
				b.WriteString(fmt.Sprintf("<span class=\"p-2 bg-gray-900 rounded-lg\">%s</span>", item.Name))
			}
		}
		b.WriteString("</span>")
		user.AddMessage(b.String())
	}
}

func UserEquip(splitMsg []string, user *gameobjects.User) {
	if len(splitMsg) < 2 {
		user.AddMessage("Usage: equip <item>")
		return
	}
	itemName := strings.Join(splitMsg[1:], " ")
	item, notFound := findItem(user, itemName)
	if notFound {
		user.AddMessage("Item not found.")
		return
	}
	if slices.Contains(user.Equipped, item) {
		user.AddMessage("You are already wearing that item.")
		return
	}
	for _, equippedItem := range user.Equipped {
		if equippedItem.Slot == item.Slot {
			user.AddMessage(fmt.Sprintf("You are already wearing the following item: %s, which is in the same item slot.", equippedItem.Name))
			return
		}
	}
	user.EquipItem(item)
}

func findItem(user *gameobjects.User, itemName string) (*gameobjects.Item, bool) {
	var item *gameobjects.Item
	for _, i := range user.Items {
		if i.Name == itemName {
			item = i
			break
		}
	}
	if item == nil {
		// search tags
		for _, i := range user.Items {
			if slices.Contains(i.Tags, itemName) {
				item = i
				break
			}
		}
	}
	if item == nil {
		return nil, true
	}
	return item, false
}

func UserUnequip(splitMsg []string, user *gameobjects.User) {
	if len(splitMsg) < 2 {
		user.AddMessage("Usage: unequip <item>")
		return
	}
	itemName := strings.Join(splitMsg[1:], " ")
	item, shouldReturn := findItem(user, itemName)
	if shouldReturn {
		return
	}
	if !slices.Contains(user.Equipped, item) {
		user.AddMessage("You are not wearing that item.")
		return
	}
	user.UnequipItem(item)
}
