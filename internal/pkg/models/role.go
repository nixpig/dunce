package models

type RoleName string

const (
	ReaderRole RoleName = "reader"
	AuthorRole RoleName = "author"
	AdminRole  RoleName = "admin"
)
