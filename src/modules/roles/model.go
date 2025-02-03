package roles

import (
	"time"
)

type Role struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_At  time.Time `json:"created_at"`
	Created_By  string    `json:"created_by"`
	Modified_At time.Time `json:"modified_at"`
	Modified_By string    `json:"modified_by"`
}
