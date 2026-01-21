package model

type User struct {
	ID       uint
	UserID   string
	Mail     string
	Password string
}

func NewUser(userID, mail, password string) *User {
	return &User{
		UserID:   userID,
		Mail:     mail,
		Password: password,
	}
}
