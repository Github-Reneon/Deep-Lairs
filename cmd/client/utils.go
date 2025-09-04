package main

import (
	"deep_lairs/internal/dbo"
	"deep_lairs/internal/gameobjects"
	"log"
	"time"
)

type LoadedUser struct {
	Ttl time.Time
	gameobjects.User
}

var LoadedUsers []LoadedUser = []LoadedUser{}

func LoadedUserMemPruner() {
	for {
		time.Sleep(time.Millisecond * 100)
		for i, loadedUser := range LoadedUsers {
			if time.Now().Before(loadedUser.Ttl) {
				continue
			}
			// remove the expired user
			log.Println("Pruning user from memory:", loadedUser.Username)
			LoadedUsers = append(LoadedUsers[:i], LoadedUsers[i+1:]...)
		}
	}
}

func FindUserMem(userName string) bool {
	for _, user := range LoadedUsers {
		if user.Username == userName {
			return true
		}
	}
	return false
}

func GetUserInMem(userName string) bool {
	user, err := dbo.LoadUser(userName)
	if err != nil {
		return false
	}
	LoadedUsers = append(LoadedUsers, LoadedUser{
		Ttl:  time.Now().Add(time.Hour * 24),
		User: user,
	})
	return true
}
