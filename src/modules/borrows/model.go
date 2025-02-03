package borrows

import "time"

type Borrow struct {
	Id              string     `json:"id"`
	User_Id         string     `json:"user_id"`
	Books           []string   `json:"books"`
	Borrowed_Time   *time.Time `json:"borrowed_time"`
	Return_Deadline *time.Time `json:"return_deadline"`
	Returned_Time   *time.Time `json:"returned_time"`
	Status          string     `json:"status"`
	Created_By      string     `json:"created_by"`
}

// user borrow
// user return
// if return late, create a new penalty record
// if not proceed just to change the status and return time
// borrow edited only when return, otherwise delete the borrow data, because it's risky when editing, since it involve book stock quantity

// user max borrow 3 books
