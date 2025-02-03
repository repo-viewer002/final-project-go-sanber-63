package utils

import "final-project/src/commons"

func IsValidRole(role string) bool {
	return role == commons.Roles.Admin || role == commons.Roles.Librarian || role == commons.Roles.Member
}
