package models

// User holds user data.
type User struct {
	ID       string
	Username string
	Password string // hashed password
}
