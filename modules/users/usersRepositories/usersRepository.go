package usersRepositories

import (
	"go_learn_project_rest_api/modules/users"
	"go_learn_project_rest_api/modules/users/usersPatterns"

	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface {
	InsertUser(*users.UserRegisterReq, bool) (*users.UserPassport, error)
}

type usersrepository struct {
	db *sqlx.DB
}

func UsersRepository(db *sqlx.DB) IUsersRepository {
	return &usersrepository{
		db: db,
	}
}

func (u *usersrepository) InsertUser(request *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	result := usersPatterns.InsertUser(u.db, request, isAdmin)

	var err error
	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	user, err := result.Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}
