package genres

import (
	"database/sql"
	"final-project/src/configs/database"
	"fmt"
)

type Repository interface {
	CreateGenreRepository(genre Genre) (Genre, error)
	GetAllGenreRepository(name string) ([]Genre, error)
	GetGenreByIdRepository(id string) (Genre, error)
	GetGenreIdByNameRepository(name string) (string, error)
	UpdateGenreByIdRepository(id string, genre Genre) (Genre, error)
	DeleteGenreByIdRepository(id string) (Genre, error)
}

type genreRepository struct{}

func NewRepository() Repository {
	return &genreRepository{}
}

func (repository *genreRepository) CreateGenreRepository(genre Genre) (Genre, error) {
	query := `
		INSERT INTO genres
		(
			name,
			description,
			created_by, 
			modified_by
		)
		VALUES
		($1, $2, $3, $4)
		RETURNING *
	`

	err := database.DB.QueryRow(query, genre.Name, genre.Description, genre.Created_By, genre.Modified_By).
		Scan(&genre.Id, &genre.Name, &genre.Description, &genre.Created_At, &genre.Created_By, &genre.Modified_At, &genre.Modified_By)

	if err != nil {
		return Genre{}, err
	}

	return genre, err
}

func (repository *genreRepository) GetAllGenreRepository(name string) ([]Genre, error) {
	var genres []Genre

	// Start building the query
	query := "SELECT * FROM genres"
	var args []interface{}

	// Check if the name parameter is provided
	if name != "" {
		// Add a WHERE clause for name filtering
		query += " WHERE name ILIKE $1"
		args = append(args, "%"+name+"%") // Using ILIKE for case-insensitive search
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return []Genre{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var genre Genre

		err = rows.Scan(&genre.Id, &genre.Name, &genre.Description, &genre.Created_At, &genre.Created_By, &genre.Modified_At, &genre.Modified_By)

		if err != nil {
			return []Genre{}, err
		}

		genres = append(genres, genre)
	}

	return genres, nil
}

func (repository *genreRepository) GetGenreByIdRepository(id string) (Genre, error) {
	var genre Genre

	query := `
		SELECT * FROM genres 
		WHERE id = $1
	`

	err := database.DB.QueryRow(query, id).
		Scan(&genre.Id, &genre.Name, &genre.Description, &genre.Created_At, &genre.Created_By, &genre.Modified_At, &genre.Modified_By)

	if err != nil {
		if err == sql.ErrNoRows {
			return Genre{}, fmt.Errorf("failed to get genre data, genre with id \"%s\" not found", id)
		}

		return Genre{}, err
	}

	return genre, nil
}

func (repository *genreRepository) GetGenreIdByNameRepository(name string) (string, error) {
	var genre Genre

	query := `
		SELECT id FROM genres 
		WHERE name = $1
	`

	err := database.DB.QueryRow(query, name).
		Scan(&genre.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("failed to get genre data, genre with name \"%s\" not found", name)
		}

		return "", err
	}

	return genre.Id, nil
}

func (repository *genreRepository) UpdateGenreByIdRepository(id string, genre Genre) (Genre, error) {
	query := `
		UPDATE genres 
		SET 
			name = $2,
			description = $3,
			modified_by = $4 
		WHERE id = $1 
		RETURNING *
	`

	err := database.DB.QueryRow(query, id, genre.Name, genre.Description, genre.Modified_By).
		Scan(&genre.Id, &genre.Name, &genre.Description, &genre.Created_At, &genre.Created_By, &genre.Modified_At, &genre.Modified_By)

	if err != nil {
		if err == sql.ErrNoRows {
			return genre, fmt.Errorf("failed updating genre, genre with id \"%s\" not found", id)
		}

		return Genre{}, err
	}

	return genre, nil
}

func (repository *genreRepository) DeleteGenreByIdRepository(id string) (Genre, error) {
	var deletedGenre Genre

	query := `
		DELETE FROM genres 
		WHERE id = $1 
		RETURNING *
	`

	err := database.DB.QueryRow(query, id).
		Scan(&deletedGenre.Id, &deletedGenre.Name, &deletedGenre.Description, &deletedGenre.Created_At, &deletedGenre.Created_By, &deletedGenre.Modified_At, &deletedGenre.Modified_By)

	if err != nil {
		if err == sql.ErrNoRows {
			return deletedGenre, fmt.Errorf("failed deleting genre, genre with id \"%s\" not found", id)
		}

		return Genre{}, err
	}

	return deletedGenre, nil
}
