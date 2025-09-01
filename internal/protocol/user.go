package protocol

type User struct {
	ID           string   `gorm:"id"`
	Username     string   `gorm:"username"`
	Password     string   `gorm:"password"`
	Email        string   `gorm:"email"`
	CharacterIds []string `gorm:"characters_ids"`
}
