package usersUsecases

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/users"
	"go_learn_project_rest_api/modules/users/usersRepositories"
)

type IUsersUsecases interface {
	InsertCustomer(*users.UserRegisterReq) (*users.UserPassport, error)
}

type usersUsecases struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecases(cfg config.IConfig, usersRepository usersRepositories.IUsersRepository) IUsersUsecases {
	return &usersUsecases{
		usersRepository: usersRepository,
		cfg:             cfg,
	}
}

func (u *usersUsecases) InsertCustomer(request *users.UserRegisterReq) (*users.UserPassport, error) {
	if err := request.BcryptHashing(); err != nil {
		return nil, err
	}

	result, err := u.usersRepository.InsertUser(request, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}
