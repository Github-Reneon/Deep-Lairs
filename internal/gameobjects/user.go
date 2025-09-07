package gameobjects

import "time"

type User struct {
	ID                 string `json:"id"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	Email              string `json:"email"`
	CurrentCharacterId string `json:"-"`
	DateCreated        time.Time
}
