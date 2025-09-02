package gameobjects

type User struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password" gorm:""`
	Email    string `json:"email" gorm:"unique"`
}
