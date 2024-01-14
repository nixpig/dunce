package models

type Site struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Owner       User   `json:"owner"`
}
