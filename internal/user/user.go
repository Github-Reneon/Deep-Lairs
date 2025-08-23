package user

type User struct {
	ID           string
	Name         string
	MessageQueue []string
}

func (u *User) GetName() string {
	if u.Name == "" {
		return u.ID
	}
	return u.Name
}
