package main

import (
	"deep_lairs/internal/dbo"
	"deep_lairs/internal/gameobjects"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoadedUser struct {
	Ttl time.Time
	gameobjects.User
}

var LoadedUsers []LoadedUser = []LoadedUser{}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	log.Printf("Genned password hash %s.\n", hash)

	return string(hash)
}

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

func PutUserInMem(userName string) bool {
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

// gets a copy not the real deal
func GetUserInMemFromName(userName string) (gameobjects.User, error) {
	for _, user := range LoadedUsers {
		if user.Username == userName {
			return user.User, nil
		}
	}
	return gameobjects.User{}, errors.New("Cannot find the user")
}
