package models

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Link     string `json:"link"` // e.g. twitter.com/nixpig
	RoleId   int    `json:"role_id"`
}
