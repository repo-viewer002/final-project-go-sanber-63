package genres

type Service interface {
	CreateGenreService(genre Genre) (Genre, error)
	GetAllGenreService(name string) ([]Genre, error)
	GetGenreByIdService(genreId string) (Genre, error)
	GetGenreIdByNameRepository(name string) (string, error)
	UpdateGenreByIdService(genreId string, genre Genre) (Genre, error)
	DeleteGenreByIdService(genreId string) (Genre, error)
}

type genreService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &genreService{
		repository,
	}
}

func (service *genreService) CreateGenreService(genre Genre) (Genre, error) {
	createdGenre, err := service.repository.CreateGenreRepository(genre)

	if err != nil {
		return Genre{}, err
	}

	return createdGenre, nil
}

func (service *genreService) GetAllGenreService(name string) ([]Genre, error) {
	genre, err := service.repository.GetAllGenreRepository(name)

	if err != nil {
		return []Genre{}, err
	}

	return genre, nil
}

func (service *genreService) GetGenreByIdService(genreId string) (Genre, error) {
	genre, err := service.repository.GetGenreByIdRepository(genreId)

	if err != nil {
		return Genre{}, err
	}

	return genre, nil
}

func (service *genreService) GetGenreIdByNameRepository(name string) (string, error) {
	genre, err := service.repository.GetGenreIdByNameRepository(name)

	if err != nil {
		return "", err
	}

	return genre, nil
}

func (service *genreService) UpdateGenreByIdService(genreId string, genre Genre) (Genre, error) {
	updatedGenre, err := service.repository.UpdateGenreByIdRepository(genreId, genre)

	if err != nil {
		return Genre{}, err
	}

	return updatedGenre, err
}

func (service *genreService) DeleteGenreByIdService(genreId string) (Genre, error) {
	deletedGenre, err := service.repository.DeleteGenreByIdRepository(genreId)

	if err != nil {
		return Genre{}, err
	}

	return deletedGenre, err
}
