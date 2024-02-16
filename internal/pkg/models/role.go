package models

import "fmt"

type RoleName string

const (
	ReaderRole RoleName = "reader"
	AuthorRole RoleName = "author"
	AdminRole  RoleName = "admin"
)

var RoleNames = []string{
	"reader", "author", "admin",
}

func (r RoleName) String() string {
	return string(r)
}

func ParseRoleName(s string) (r RoleName, e error) {
	roleNames := map[string]RoleName{
		"reader": ReaderRole,
		"author": AuthorRole,
		"admin":  AdminRole,
	}

	roleName := RoleName(s)

	_, ok := roleNames[s]
	if !ok {
		return r, fmt.Errorf("cannot parse '%s' as 'RoleName'", s)
	}

	return roleName, nil
}
