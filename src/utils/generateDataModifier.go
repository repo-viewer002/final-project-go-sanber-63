package utils

import "final-project/src/commons"

func GenerateDataModifier(role string, username string, modifier *string) {
	switch role {
	case commons.Roles.Admin:
		*modifier = commons.Roles.Admin + " " + username
	case commons.Roles.Librarian:
		*modifier = commons.Roles.Librarian + " " + username
	case commons.Roles.Member:
		*modifier = "system"
	default:
		*modifier = "system"
	}
}
