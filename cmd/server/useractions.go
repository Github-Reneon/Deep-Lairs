package main

import (
	"deep_lairs/internal/gameobjects"
	"deep_lairs/internal/protocol"
	"fmt"
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
		user.AddMessage(fmt.Sprintf(protocol.SAY, user.GetName(), strings.ToUpper(splitMsg[1][:1])+strings.Join(splitMsg[1:], " ")[1:]))
		for _, u := range user.Location.Users {
			if u != user {
				u.AddMessage(fmt.Sprintf(protocol.SAY, user.GetName(), strings.Join(splitMsg[1:], " ")))
			}
		}
	}
}

func UserLook(splitMsg []string, user *gameobjects.User) {
	if len(splitMsg) > 1 {
		user.AddMessage(fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, "look "+strings.Join(splitMsg[1:], " ")))
	} else {
		user.AddMessage(fmt.Sprintf(protocol.LOOK, user.Location.QuickLook, user.Location.Look, user.Location.LookImage))
		users := []string{}
		for _, foundUser := range user.Location.Users {
			if foundUser != user {
				users = append(users, foundUser.GetName())
			}
		}
		if len(users) > 0 {
			if len(users) >= 10 {
				user.AddMessage(fmt.Sprintf("You see many adventurers here. %d in total.", len(users)))
			}
			for _, u := range users {
				user.AddMessage(fmt.Sprintf("You see %s here.", u))
			}
		} else {
			user.AddMessage("You don't see any other adventurers here.")
		}
	}
}
