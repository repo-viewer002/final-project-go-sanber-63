package auth

import (
	"errors"
	"final-project/src/commons/middlewares"
	"final-project/src/utils"
)

type Service interface {
	LoginService(credentials Credentials) (string, string, error)
}

type authService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &authService{
		repository,
	}
}

func (service *authService) LoginService(credentials Credentials) (string, string, error) {
	validUser, err := service.repository.ValidateUsernameAndEmail(credentials.Identifier)

	if err != nil {
		return "", "", err
	}

	if validPassword := utils.CompareWithHash(credentials.Password, validUser.Password); !validPassword {
		return "", "", errors.New("invalid credentials")
	}

	token, err := middlewares.CreateToken(validUser.Id, validUser.Username, validUser.Email, validUser.Role)

	if err != nil {
		return "", "", err
	}

	return token, validUser.Role, nil
}
