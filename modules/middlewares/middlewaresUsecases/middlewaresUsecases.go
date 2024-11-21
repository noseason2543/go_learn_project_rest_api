package middlewaresUsecases

import middlewaresrepository "go_learn_project_rest_api/modules/middlewares/middlewaresRepository"

type IMiddlewaresUsecases interface {
}

type middlewaresUsecases struct {
	middlewaresRepository middlewaresrepository.IMiddlewaresRepository
}

func MiddlewaresUsecases(m middlewaresrepository.IMiddlewaresRepository) IMiddlewaresUsecases {
	return &middlewaresUsecases{
		middlewaresRepository: m,
	}
}
