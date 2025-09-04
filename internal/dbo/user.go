package dbo

import (
	"deep_lairs/internal/gameobjects"
)

func LoadUser(username string) (gameobjects.User, error) {
	var user gameobjects.User
	err := db.QueryRow("SELECT id, username, password, email FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return gameobjects.User{}, err
	}
	return user, nil
}

func CreateUser(username, password, email string) error {
	_, err := db.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?, datetime('now'))", username, password, email)
	return err
}
