package protocol

type User struct {
	ID           string   `json:"id"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	Email        string   `json:"email"`
	CharacterIds []string `json:"characters_ids"`
}
