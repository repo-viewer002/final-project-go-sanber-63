package borrows

import (
	"database/sql"
	"final-project/src/commons"
	"final-project/src/configs/database"
	"fmt"
	"time"
)

type Repository interface {
	BorrowBookRepository(borrow Borrow) (Borrow, error)
	ReturnBookRepository(borrowId string) (Borrow, error)
	// GetAllBorrowRepository(searchType string, genres ...string) ([]Book, error)
	// GetBorrowByIdRepository(borrowId string) (Book, error)
	// DeleteBorrowRepository(searchBook SearchBook) ([]Book, error)
}

type borrowRepository struct{}

func NewRepository() Repository {
	return &borrowRepository{}
}

func (repository *borrowRepository) BorrowBookRepository(borrow Borrow) (Borrow, error) {

	// check user penalized status, is it more than current time
	// if yes return error of user is penalized
	// if not, clear penalty_duration and change user status to active

	repository.CheckUserStatusAndPenaltyDuration(borrow.User_Id)
	repository.CheckUserTotalBorrowed(borrow.User_Id)

	var bookNames []string

	tx, err := database.DB.Begin()
	if err != nil {
		return Borrow{}, err
	}

	query := `
		INSERT INTO borrows 
		(
			user_id, 
			return_deadline,
			created_by
		)
		VALUES 
		(
			$1, 
			CURRENT_TIMESTAMP + INTERVAL '7 days',
			$2
		)
		RETURNING 
			id, 
			user_id, 
			borrowed_time, 
			return_deadline, 
			returned_time, 
			status,
			created_by
	`

	err = tx.QueryRow(query, borrow.User_Id, borrow.Created_By).
		Scan(&borrow.Id, &borrow.User_Id, &borrow.Borrowed_Time, &borrow.Return_Deadline, &borrow.Returned_Time, &borrow.Status, &borrow.Created_By)

	if err != nil {
		tx.Rollback()
		return Borrow{}, err
	}

	borrowedBooksQuery := `
		INSERT INTO borrowed_books 
		(
			borrow_id, 
			book_id
		)
		VALUES ($1, $2)
		RETURNING 
			(SELECT name FROM books WHERE id = $2)
	`
	for _, bookId := range borrow.Books {
		duplicated, err := repository.CheckUserDuplicatedBookBorrowed(borrow.User_Id, bookId)

		if err != nil {
			tx.Rollback()
			return Borrow{}, fmt.Errorf("failed to check if book with id \"%s\" is already borrowed: %w", bookId, err)
		}

		if duplicated {
			tx.Rollback()
			return Borrow{}, fmt.Errorf("user with id \"%s\" has already borrowed the book with id \"%s\"", borrow.User_Id, bookId)
		}

		var bookName string

		err = tx.QueryRow(borrowedBooksQuery, borrow.Id, bookId).
			Scan(&bookName)

		if err != nil {
			tx.Rollback()

			return Borrow{}, err
		}

		err = repository.DecreaseBookStockAndIncreaseBorrow(bookId)

		if err != nil {
			tx.Rollback()

			return Borrow{}, err
		}

		bookNames = append(bookNames, bookName)
	}

	err = tx.Commit()

	if err != nil {
		return Borrow{}, err
	}

	borrow.Books = bookNames

	return borrow, nil
}

func (repository *borrowRepository) ReturnBookRepository(borrowId string) (Borrow, error) {
	tx, err := database.DB.Begin()

	if err != nil {
		return Borrow{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	borrowStatus, err := repository.CheckBorrowStatus(borrowId)
	if err != nil {
		tx.Rollback()

		return Borrow{}, err
	}

	if borrowStatus != "borrowed" {
		tx.Rollback()

		return Borrow{}, fmt.Errorf("user has already returned this book")
	}

	late, err := repository.CheckExceedingReturnDeadline(borrowId)
	if err != nil {
		tx.Rollback()

		return Borrow{}, err
	}

	newStatus := "returned"

	if late {
		newStatus = "overdue"

		var userId string
		var overdueDays int

		getUserQuery :=
			`
			SELECT 
				user_id,
				EXTRACT(DAY FROM CURRENT_TIMESTAMP - return_deadline)
			FROM 
				borrows 
			WHERE 
				id = $1
			`

		err = tx.QueryRow(getUserQuery, borrowId).Scan(&userId, &overdueDays)

		if err != nil {
			tx.Rollback()

			return Borrow{}, err
		}

		totalPenalty := overdueDays * commons.PENALTY_AMOUNT_PER_DAY

		penaltyQuery :=
			`
			INSERT INTO penalties 
			(
				borrow_id,
				total_amount, 
			)
			VALUES 
			(
				$1, 
				$2, 
			)
		`
		_, err = tx.Exec(penaltyQuery, borrowId, totalPenalty) // i want total amount to be calculated of how many days late
		if err != nil {
			tx.Rollback()
			return Borrow{}, err
		}

		// Penalize user
		penalizeUserQuery :=
			`
			UPDATE 
				users
			SET 
				is_penalized = TRUE, 
				penalty_duration = CURRENT_TIMESTAMP + INTERVAL '3 days', 
				status = 'penalized'
			WHERE user_id = $1
		`
		_, err = tx.Exec(penalizeUserQuery, userId)
		if err != nil {
			tx.Rollback()
			return Borrow{}, err
		}
	}

	// Update borrow record to mark as returned
	var returnedBook Borrow
	updateQuery :=
		`
		UPDATE 
			borrows
		SET 
			status = $2, 
			returned_time = CURRENT_TIMESTAMP
		WHERE 
			id = $1
		RETURNING *
	`
	err = tx.QueryRow(updateQuery, borrowId, newStatus).
		Scan(&returnedBook.Id, &returnedBook.User_Id, &returnedBook.Borrowed_Time, &returnedBook.Return_Deadline, &returnedBook.Returned_Time, &returnedBook.Status, &returnedBook.Created_By)

	if err != nil {
		tx.Rollback()
		return Borrow{}, err
	}

	bookNames, err := repository.IncreaseBookStock(borrowId)
	if err != nil {
		return Borrow{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Borrow{}, err
	}

	returnedBook.Books = bookNames

	return returnedBook, nil
}

func (repostitory *borrowRepository) CheckExceedingReturnDeadline(borrowId string) (bool, error) {
	var returnDeadline time.Time
	query :=
		`
		SELECT 
			return_deadline
		FROM
			borrows
		WHERE
			id = $1
	`

	err := database.DB.QueryRow(query, borrowId).
		Scan(&returnDeadline)

	if err != nil {
		return false, err
	}

	currentTime := time.Now()

	if currentTime.After(returnDeadline) {
		return true, nil
	}

	return false, nil
}

func (repostitory *borrowRepository) CheckBorrowStatus(borrowId string) (string, error) {
	var borrowStatus string

	query :=
		`
		SELECT 
			status
		FROM
			borrows
		WHERE
			id = $1

	`

	err := database.DB.QueryRow(query, borrowId).
		Scan(&borrowStatus)

	if err != nil {
		return "", err
	}

	return borrowStatus, nil
}

func (repository *borrowRepository) CheckUserStatusAndPenaltyDuration(userId string) error {
	var isPenalized bool
	var penaltyDuration time.Time
	var status string

	query :=
		`
		SELECT
			is_penalized,
			penalty_duration,
			status
		FROM users
		WHERE
			user_id = $1
	`

	err := database.DB.QueryRow(query, userId).
		Scan(&isPenalized, &penaltyDuration, &status)

	if err != nil {
		return err
	}

	currentTime := time.Now()

	if isPenalized && penaltyDuration.After(currentTime) {
		return fmt.Errorf("failed borrow books, user with id %s status is %s, with penalty duration until %s", userId, status, penaltyDuration)
	}

	if isPenalized && penaltyDuration.Before(currentTime) {
		updateQuery :=
			`
			UPDATE 
				users
			SET 
				is_penalized = FALSE, 
				penalty_duration = NULL, 
				status = 'active'
			WHERE 
				user_id = $1
		`

		_, err := database.DB.Exec(updateQuery, userId)

		if err != nil {
			return fmt.Errorf("failed to update user status after penalty expiration: %w", err)
		}
	}

	return nil
}

func (repository *borrowRepository) DecreaseBookStockAndIncreaseBorrow(bookId string) error {
	var stock int

	checkStockQuery := `SELECT stock FROM books WHERE id = $1`

	err := database.DB.QueryRow(checkStockQuery, bookId).Scan(&stock)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("book with id \"%s\" not found", bookId)
		}

		return fmt.Errorf("failed to get stock for book with id \"%s\": %w", bookId, err)
	}

	if stock <= 0 {
		return fmt.Errorf("insufficient stock for book with id \"%s\", stock is 0 or less", bookId)
	}

	query := `
		UPDATE 
			books
		SET
			stock = stock - 1,
			borrowed = borrowed + 1
		WHERE
			id = $1
	`

	result, err := database.DB.Exec(query, bookId)
	if err != nil {
		return fmt.Errorf("failed to update stock for book with id \"%s\": %w", bookId, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update stock for book with id \"%s\": %w", bookId, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("book with id \"%s\" not found", bookId)
	}

	return nil
}

func (repository *borrowRepository) IncreaseBookStock(borrowId string) ([]string, error) {
	var bookNames []string
	// Query to get the book IDs from borrowed_books
	getBookIdsQuery := `
		SELECT
			book_id
		FROM
			borrowed_books
		WHERE
			borrow_id = $1
	`

	// Directly execute the query without a transaction context
	rows, err := database.DB.Query(getBookIdsQuery, borrowId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through all the book ids
	for rows.Next() {
		var bookId string

		err = rows.Scan(&bookId)
		if err != nil {
			return nil, err
		}

		fmt.Println("Book Id:", bookId)

		// Query to get the current stock and borrowed values for the book
		checkStockQuery := `SELECT name, stock, borrowed FROM books WHERE id = $1`

		var stock, borrowed int
		var bookName string

		err = database.DB.QueryRow(checkStockQuery, bookId).Scan(&bookName,&stock, &borrowed)
		fmt.Println("err:", stock)
		fmt.Println("Stock:", stock)
		fmt.Println("Borrowed:", borrowed)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("book with id \"%s\" not found", bookId)
			}
			return nil, fmt.Errorf("failed to get stock for book with id \"%s\": %w", bookId, err)
		}

		// Check if borrowed is greater than 0, and then update stock and borrowed count
		if borrowed <= 0 {
			return nil, fmt.Errorf("no borrowed books found for book with id \"%s\"", bookId)
		}

		// Update the book stock and borrowed count
		updateQuery := `
		UPDATE 
			books
		SET
			stock = stock + 1,
			borrowed = borrowed - 1
		WHERE
			id = $1
		`

		_, err = database.DB.Exec(updateQuery, bookId)
		if err != nil {
			return nil, fmt.Errorf("failed to update stock for book with id \"%s\": %w", bookId, err)
		}

		bookNames = append(bookNames, bookName)
	}

	return bookNames, nil
}


func (repository *borrowRepository) CheckUserTotalBorrowed(userId string) error {
	var userTotalBorrowed int

	checkBorrowedCountQuery :=
		`
		SELECT 
			COUNT(*) 
		FROM 
			borrowed_books 
		WHERE 
			borrow_id IN (SELECT id FROM borrows WHERE user_id = $1 AND status != 'returned')
		`
	err := database.DB.QueryRow(checkBorrowedCountQuery, userId).Scan(&userTotalBorrowed)

	if err != nil {
		return fmt.Errorf("failed to check borrowed count for user with id \"%s\": %w", userId, err)
	}

	if userTotalBorrowed >= 3 {
		return fmt.Errorf("user with id \"%s\" has already borrowed 3 books", userId)
	}

	return nil
}

func (repository *borrowRepository) CheckUserDuplicatedBookBorrowed(userId string, bookId string) (bool, error) {
	var duplicatedBorrowedBook int

	checkExistingBookQuery :=
		`
		SELECT 
			COUNT(*) 
		FROM 
			borrowed_books 
		WHERE 
			borrow_id IN (SELECT id FROM borrows WHERE user_id = $1 AND
			status != 'returned') AND
			book_id = $2
		`
	err := database.DB.QueryRow(checkExistingBookQuery, userId, bookId).Scan(&duplicatedBorrowedBook)

	if err != nil {
		return false, err
	}

	return duplicatedBorrowedBook > 0, nil
}
