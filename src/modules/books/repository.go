package books

import (
	"database/sql"
	"errors"
	"final-project/src/configs/database"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type Repository interface {
	CreateBookRepository(book Book) (Book, error)
	GetAllBookRepository(searchBook SearchBook) ([]Book, error)
	GetAllBookByGenreRepository(searchType string, genres ...string) ([]Book, error)
	GetBookByIdRepository(bookId string) (Book, error)
	UpdateBookByIdRepository(bookId string, book Book) (Book, error)
	DeleteBookByIdRepository(bookId string) (Book, error)
}

type bookRepository struct{}

func NewRepository() Repository {
	return &bookRepository{}
}

func (repository *bookRepository) CreateBookRepository(book Book) (Book, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return Book{}, fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
      INSERT INTO books (
          name, 
          description, 
          authors, 
          publisher, 
          publish_year, 
          stock, 
          created_by, 
          modified_by
      ) 
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
      RETURNING 
          id, 
          name, 
          description, 
          authors, 
          publisher, 
          publish_year, 
          stock, 
          borrowed,
          created_at, 
          created_by, 
          modified_at, 
          modified_by`

	var result Book

	err = tx.QueryRow(
		query,
		book.Name,
		book.Description,
		book.Authors,
		book.Publisher,
		book.Publish_Year,
		book.Stock,
		book.Created_By,
		book.Modified_By,
	).Scan(
		&result.Id,
		&result.Name,
		&result.Description,
		&result.Authors,
		&result.Publisher,
		&result.Publish_Year,
		&result.Stock,
		&result.Borrowed,
		&result.Created_At,
		&result.Created_By,
		&result.Modified_At,
		&result.Modified_By,
	)

	if err != nil {
		return Book{}, fmt.Errorf("failed to insert and scan book: %v", err)
	}

	for _, genreName := range book.Genres {
		var genreId string
		err := tx.QueryRow("SELECT id FROM genres WHERE name = $1", genreName).Scan(&genreId)
		if err != nil {
			return Book{}, fmt.Errorf("genre %s does not exist", genreName)
		}

		_, err = tx.Exec(
			"INSERT INTO book_genres (book_id, genre_id) VALUES ($1, $2)",
			result.Id,
			genreId,
		)
		if err != nil {
			return Book{}, fmt.Errorf("failed to insert book_genre: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return Book{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	result.Genres = book.Genres

	return result, nil
}

func (repository *bookRepository) GetAllBookRepository(searchBook SearchBook) ([]Book, error) {
	var books []Book
	var args []interface{}
	argPosition := 1

	// Base query with genre aggregation
	baseQuery := `
			WITH filtered_books AS (
					SELECT DISTINCT b.id
					FROM books b
					WHERE 1=1
	`

	// Add basic filters
	if searchBook.Name != "" {
		baseQuery += fmt.Sprintf(" AND b.name ILIKE $%d", argPosition)
		args = append(args, "%"+searchBook.Name+"%")
		argPosition++
	}

	if searchBook.Authors != "" {
		baseQuery += fmt.Sprintf(" AND b.authors ILIKE $%d", argPosition)
		args = append(args, "%"+searchBook.Authors+"%")
		argPosition++
	}

	if searchBook.Publisher != "" {
		baseQuery += fmt.Sprintf(" AND b.publisher ILIKE $%d", argPosition)
		args = append(args, "%"+searchBook.Publisher+"%")
		argPosition++
	}

	if searchBook.Publish_Year != "" {
		baseQuery += fmt.Sprintf(" AND b.publish_year = $%d", argPosition)
		args = append(args, searchBook.Publish_Year)
		argPosition++
	}

	// Handle genre filtering if genres are provided
	if len(searchBook.Genres) > 0 && searchBook.Genres[0] != "" {
		// Validate genres first
		for _, genreName := range searchBook.Genres {
			var genreId string
			err := database.DB.QueryRow("SELECT id FROM genres WHERE name = $1", genreName).Scan(&genreId)
			if err != nil {
				return nil, fmt.Errorf("genre %s does not exist", genreName)
			}
		}

		// Add genre filtering based on search type
		baseQuery += `
					AND EXISTS (
							SELECT 1 FROM book_genres bg
							JOIN genres g ON g.id = bg.genre_id
							WHERE bg.book_id = b.id
							AND g.name = ANY($` + fmt.Sprintf("%d", argPosition) + `)
					)`
		args = append(args, pq.Array(searchBook.Genres))
		argPosition++

		if searchBook.Genre_Search_Type == "all" {
			baseQuery += fmt.Sprintf(`
							AND (
									SELECT COUNT(DISTINCT g2.name)
									FROM book_genres bg2
									JOIN genres g2 ON g2.id = bg2.genre_id
									WHERE bg2.book_id = b.id
									AND g2.name = ANY($%d)
							) = $%d`, argPosition-1, argPosition)
			args = append(args, len(searchBook.Genres))
			argPosition++
		}
	}

	// Close the CTE and add the main query
	mainQuery := baseQuery + `
			)
			SELECT 
					b.*,
					STRING_AGG(DISTINCT g.name, ', ' ORDER BY g.name) AS genres
			FROM filtered_books fb
			JOIN books b ON b.id = fb.id
			LEFT JOIN book_genres bg ON bg.book_id = b.id
			LEFT JOIN genres g ON g.id = bg.genre_id
			GROUP BY 
					b.id, b.name, b.description, b.authors, b.publisher,
					b.publish_year, b.stock, b.borrowed, b.created_at,
					b.created_by, b.modified_at, b.modified_by
			ORDER BY b.name
	`

	// Execute query
	rows, err := database.DB.Query(mainQuery, args...)
	if err != nil {
		return []Book{}, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		var book Book
		var genres string
		err = rows.Scan(
			&book.Id,
			&book.Name,
			&book.Description,
			&book.Authors,
			&book.Publisher,
			&book.Publish_Year,
			&book.Stock,
			&book.Borrowed,
			&book.Created_At,
			&book.Created_By,
			&book.Modified_At,
			&book.Modified_By,
			&genres,
		)
		if err != nil {
			return []Book{}, fmt.Errorf("failed to scan row: %v", err)
		}

		if genres != "" {
			book.Genres = strings.Split(genres, ", ")
		} else {
			book.Genres = []string{}
		}

		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return []Book{}, fmt.Errorf("error iterating rows: %v", err)
	}

	return books, nil
}

func (repository *bookRepository) GetAllBookByGenreRepository(searchType string, genres ...string) ([]Book, error) {
	var books []Book

	genreCount := len(genres)
	if genreCount == 0 {
		return nil, fmt.Errorf("no genres provided")
	}

	query := `
			SELECT name 
			FROM genres 
			WHERE name = ANY($1)
	`
	rows, err := database.DB.Query(query, pq.Array(genres))
	if err != nil {
		return nil, fmt.Errorf("failed to validate genres: %v", err)
	}
	defer rows.Close()

	validGenres := make(map[string]bool)
	for rows.Next() {
		var genre string
		if err := rows.Scan(&genre); err != nil {
			return nil, fmt.Errorf("failed to scan genre: %v", err)
		}
		validGenres[genre] = true
	}

	var invalidGenres []string
	for _, genre := range genres {
		if !validGenres[genre] {
			invalidGenres = append(invalidGenres, genre)
		}
	}

	if len(invalidGenres) > 0 {
		return nil, fmt.Errorf("invalid genres provided: %s", strings.Join(invalidGenres, ", "))
	}

	placeholders := make([]string, genreCount)
	args := make([]interface{}, genreCount)
	for i, genre := range genres {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = genre
	}

	var groupByAndHaving string
	if searchType == "all" {
		groupByAndHaving = fmt.Sprintf("HAVING COUNT(DISTINCT g1.name) = %d", genreCount)
	} else if searchType == "any" {
		groupByAndHaving = ""
	} else {
		return nil, errors.New("invalid search type, please choose either \"any\" (search book based on any matching genres) or \"all\" (search book based on all matching genres)")
	}

	mainQuery := fmt.Sprintf(`
			WITH matching_books AS (
					SELECT DISTINCT b.id
					FROM books b
					JOIN book_genres bg1 ON bg1.book_id = b.id
					JOIN genres g1 ON g1.id = bg1.genre_id
					WHERE g1.name IN (%s)
					GROUP BY b.id
					%s
			)
			SELECT 
					b.*,
					STRING_AGG(DISTINCT g2.name, ', ' ORDER BY g2.name) AS genres
			FROM matching_books mb
			JOIN books b ON b.id = mb.id
			LEFT JOIN book_genres bg2 ON bg2.book_id = b.id
			LEFT JOIN genres g2 ON g2.id = bg2.genre_id
			GROUP BY 
					b.id, b.name, b.description, b.authors, b.publisher, 
					b.publish_year, b.stock, b.borrowed, b.created_at, 
					b.created_by, b.modified_at, b.modified_by
			ORDER BY b.name;
	`, strings.Join(placeholders, ", "), groupByAndHaving)

	rows, err = database.DB.Query(mainQuery, args...)
	if err != nil {
		return []Book{}, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		var bookGenres string
		err = rows.Scan(
			&book.Id,
			&book.Name,
			&book.Description,
			&book.Authors,
			&book.Publisher,
			&book.Publish_Year,
			&book.Stock,
			&book.Borrowed,
			&book.Created_At,
			&book.Created_By,
			&book.Modified_At,
			&book.Modified_By,
			&bookGenres,
		)
		if err != nil {
			return []Book{}, fmt.Errorf("failed to scan row: %v", err)
		}

		book.Genres = strings.Split(bookGenres, ", ")
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return []Book{}, fmt.Errorf("error iterating rows: %v", err)
	}

	return books, nil
}

func (repository *bookRepository) GetBookByIdRepository(bookId string) (Book, error) {
	var book Book

	query := `
		SELECT 
    books.*,
    STRING_AGG(genres.name, ', ') AS genres
		FROM 
			books 
		LEFT JOIN 
			book_genres ON book_genres.book_id = books.id 
		LEFT JOIN 
			genres ON genres.id = book_genres.genre_id
		WHERE 
			books.id = $1
		GROUP BY 
			books.id;
	`

	var genres string

	err := database.DB.QueryRow(query, bookId).
		Scan(&book.Id, &book.Name, &book.Description, &book.Authors, &book.Publisher, &book.Publish_Year, &book.Stock, &book.Borrowed, &book.Created_At, &book.Created_By, &book.Modified_At, &book.Modified_By, &genres)

	if err != nil {
		if err == sql.ErrNoRows {
			return Book{}, fmt.Errorf("failed to get book data, book with id \"%s\" not found", bookId)
		}

		return Book{}, err
	}

	book.Genres = strings.Split(genres, ", ")

	return book, nil
}

func (repository *bookRepository) UpdateBookByIdRepository(bookId string, book Book) (Book, error) {
	bookGenres := strings.Join(book.Genres, ", ")

	tx, err := database.DB.Begin()
	if err != nil {
		return Book{}, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	updateQuery := `
		UPDATE books 
		SET 
			name = COALESCE(NULLIF($2, ''), name),
			description = COALESCE(NULLIF($3, ''), description),
			authors = COALESCE(NULLIF($4, ''), authors), 
			publisher = COALESCE(NULLIF($5, ''), publisher),
			publish_year = COALESCE(NULLIF($6, 0), publish_year),
			stock = COALESCE(NULLIF($7, 0), stock),
			borrowed = COALESCE(NULLIF($8, 0), borrowed),
			modified_by = COALESCE(NULLIF($9, ''), modified_by)
		WHERE id = $1
		RETURNING id, name, description, authors, publisher, publish_year, stock, borrowed, created_at, created_by, modified_at, modified_by
	`

	var updatedBook Book
	err = tx.QueryRow(updateQuery, bookId, book.Name, book.Description, book.Authors, book.Publisher, book.Publish_Year, book.Stock, book.Borrowed, book.Modified_By).
		Scan(&updatedBook.Id, &updatedBook.Name, &updatedBook.Description, &updatedBook.Authors, &updatedBook.Publisher, &updatedBook.Publish_Year, &updatedBook.Stock, &updatedBook.Borrowed, &updatedBook.Created_At, &updatedBook.Created_By, &updatedBook.Modified_At, &updatedBook.Modified_By)

	if err != nil {
		if err == sql.ErrNoRows {
			return updatedBook, fmt.Errorf("failed updating book, book with id \"%s\" not found", bookId)
		}
		return Book{}, fmt.Errorf("failed updating book: %v", err)
	}

	deleteGenresQuery := `
		DELETE FROM book_genres
		WHERE book_id = $1
	`

	_, err = tx.Exec(deleteGenresQuery, bookId)
	if err != nil {
		return Book{}, fmt.Errorf("failed deleting old genres: %v", err)
	}

	insertGenresQuery := `
		INSERT INTO book_genres (book_id, genre_id)
		SELECT $1, id 
		FROM genres 
		WHERE name = ANY(string_to_array($2, ', '))
		AND NOT EXISTS (
			SELECT 1 
			FROM book_genres 
			WHERE book_id = $1 
			AND genre_id = genres.id
		)
	`

	_, err = tx.Exec(insertGenresQuery, bookId, bookGenres)
	if err != nil {
		return Book{}, fmt.Errorf("failed inserting new genres: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return Book{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	updatedBook.Genres = book.Genres

	return updatedBook, nil
}

func (repository *bookRepository) DeleteBookByIdRepository(bookId string) (Book, error) {
	var deletedBook Book

	query := `
		DELETE FROM books 
		WHERE id = $1 
		RETURNING *
	`

	err := database.DB.QueryRow(query, bookId).
		Scan(&deletedBook.Id, &deletedBook.Name, &deletedBook.Description, &deletedBook.Description, &deletedBook.Authors, &deletedBook.Publisher, &deletedBook.Publish_Year, &deletedBook.Stock, &deletedBook.Created_At, &deletedBook.Created_By, &deletedBook.Modified_At, &deletedBook.Modified_By)

	if err != nil {
		if err == sql.ErrNoRows {
			return deletedBook, fmt.Errorf("failed deleting book, book with id \"%s\" not found", bookId)
		}

		return Book{}, err
	}

	return deletedBook, nil
}
