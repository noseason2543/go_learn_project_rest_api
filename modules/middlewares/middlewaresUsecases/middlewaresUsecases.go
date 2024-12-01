package middlewaresUsecases

import (
	"go_learn_project_rest_api/modules/middlewares"
	middlewaresrepository "go_learn_project_rest_api/modules/middlewares/middlewaresRepository"
)

type IMiddlewaresUsecases interface {
	FindAccessToken(userId, token string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresUsecases struct {
	middlewaresRepository middlewaresrepository.IMiddlewaresRepository
}

func MiddlewaresUsecases(m middlewaresrepository.IMiddlewaresRepository) IMiddlewaresUsecases {
	return &middlewaresUsecases{
		middlewaresRepository: m,
	}
}

func (u *middlewaresUsecases) FindAccessToken(userId, token string) bool {
	return u.middlewaresRepository.FindAccessToken(userId, token)
}

func (u *middlewaresUsecases) FindRole() ([]*middlewares.Role, error) {
	return u.middlewaresRepository.FindRole()
}
