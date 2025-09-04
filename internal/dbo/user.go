package dbo

import (
	"deep_lairs/internal/gameobjects"
	"log"
)

func LoadUser(username string) (gameobjects.User, error) {
	var user gameobjects.User
	err := db.QueryRow("SELECT id, username, password, email FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		log.Println("Could not load user from database:", username)
		return gameobjects.User{}, err
	}
	log.Println("Loaded use from database:", username)
	return user, nil
}

func LoadUserFromId(id string) (gameobjects.User, error) {
	var user gameobjects.User
	err := db.QueryRow("SELECT id, username, password, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		log.Println("Could not load user from database:", id)
		return gameobjects.User{}, err
	}
	log.Println("Loaded use from database:", user.Username)
	return user, nil
}

func CreateUser(username, password, email string) error {
	log.Println("Creating user in database:", username)
	_, err := db.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", username, password, email)
	if err != nil {
		log.Println("Could not create user in database:", username, "error:", err)
	}
	return err
}
