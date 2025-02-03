package borrows

type Service interface {
	BorrowBookService(borrow Borrow) (Borrow, error)
	ReturnBookService(borrowId string) (Borrow, error)
}

type borrowService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &borrowService{
		repository,
	}
}

func (service *borrowService) BorrowBookService(borrow Borrow) (Borrow, error) {
	borrowData, err := service.repository.BorrowBookRepository(borrow)

	if err != nil {
		return Borrow{}, err
	}

	return borrowData, nil
}

func (service *borrowService) ReturnBookService(borrowId string) (Borrow, error) {
	borrowData, err := service.repository.ReturnBookRepository(borrowId)

	if err != nil {
		return Borrow{}, err
	}

	return borrowData, nil
}
