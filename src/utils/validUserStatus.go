package utils

import "final-project/src/commons"

func IsValidStatus(status string) bool {
	return status == commons.UserStatus.Active || status == commons.UserStatus.Suspended || status == commons.UserStatus.Deactivated
}
