package books

type Service interface {
	CreateBookService(book Book) (Book, error)
	GetAllBookService(searchBook SearchBook) ([]Book, error)
	GetAllBookByGenreService(searchType string, genres ...string) ([]Book, error)
	GetBookByIdService(bookId string) (Book, error)
	UpdateBookByIdService(bookId string, book Book) (Book, error)
	DeleteBookByIdService(bookId string) (Book, error)
}

type bookService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &bookService{
		repository,
	}
}

func (service *bookService) CreateBookService(book Book) (Book, error) {
	createdBook, err := service.repository.CreateBookRepository(book)

	if err != nil {
		return Book{}, err
	}

	return createdBook, nil
}

func (service *bookService) GetAllBookService(searchBook SearchBook) ([]Book, error) {
	book, err := service.repository.GetAllBookRepository(searchBook)

	if err != nil {
		return []Book{}, err
	}

	return book, nil
}

func (service *bookService) GetAllBookByGenreService(searchType string, genres ...string) ([]Book, error) {
	books, err := service.repository.GetAllBookByGenreRepository(searchType, genres...)

	if err != nil {
		return nil, err
	}

	return books, nil
}

func (service *bookService) GetBookByIdService(bookId string) (Book, error) {
	book, err := service.repository.GetBookByIdRepository(bookId)

	if err != nil {
		return Book{}, err
	}

	return book, nil
}

func (service *bookService) UpdateBookByIdService(bookId string, book Book) (Book, error) {
	updatedBook, err := service.repository.UpdateBookByIdRepository(bookId, book)

	if err != nil {
		return Book{}, err
	}

	return updatedBook, err
}

func (service *bookService) DeleteBookByIdService(bookId string) (Book, error) {
	deletedBook, err := service.repository.DeleteBookByIdRepository(bookId)

	if err != nil {
		return Book{}, err
	}

	return deletedBook, err
}
