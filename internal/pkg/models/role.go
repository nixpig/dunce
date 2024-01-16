package models

// TODO: look to use enums for roles - how to serialise/deserialise JSON??
type Role struct {
	Id   int    `json:"id"`
	Name string `json:"name" required:"required,max=100"`
}
