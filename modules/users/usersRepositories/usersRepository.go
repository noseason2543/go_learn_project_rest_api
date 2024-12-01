package usersRepositories

import (
	"context"
	"fmt"
	"go_learn_project_rest_api/modules/users"
	"go_learn_project_rest_api/modules/users/usersPatterns"
	"time"

	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface {
	InsertUser(*users.UserRegisterReq, bool) (*users.UserPassport, error)
	FindOneUserByEmail(string) (*users.UserCredentialCheck, error)
	InsertOauth(*users.UserPassport) error
	FindOneOauth(string) (*users.Oauth, error)
	UpdateOauth(*users.UserToken) error
	GetProfile(string) (*users.User, error)
	DeleteOauth(string) error
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

func (u *usersrepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
		SELECT 
			id,
			password,
			email,
			role_id,
			username
		FROM users WHERE email = $1;
	`

	user := new(users.UserCredentialCheck)
	if err := u.db.Get(user, query, email); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (u *usersrepository) InsertOauth(user *users.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		INSERT INTO oauth (
			user_id,
			refresh_token,
			access_token
		) VALUES ($1, $2, $3) RETURNING id;
	`
	if err := u.db.QueryRowContext(
		ctx,
		query,
		user.User.Id,
		user.Token.RefreshToken,
		user.Token.AccessToken,
	).Scan(&user.Token.Id); err != nil {
		return fmt.Errorf("insert oauth error: %v", err)
	}

	return nil
}

func (u *usersrepository) FindOneOauth(refreshToken string) (*users.Oauth, error) {
	query := `
		SELECT
			id,
			user_id
		FROM oauth WHERE refresh_token = $1;
	`
	oauthUser := new(users.Oauth)
	if err := u.db.Get(oauthUser, query, refreshToken); err != nil {
		return nil, fmt.Errorf("oauth not found")
	}
	return oauthUser, nil
}

func (u *usersrepository) UpdateOauth(req *users.UserToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		UPDATE oauth SET
			access_token = :access_token,
			refresh_token = :refresh_token
		WHERE id = :id;
	`
	if _, err := u.db.NamedExecContext(
		ctx,
		query,
		req,
	); err != nil {
		return fmt.Errorf("execution update oauth not completed: %v", err)
	}

	return nil
}

func (u *usersrepository) GetProfile(userId string) (*users.User, error) {
	query := `
		SELECT 
			id,
			email,
			username,
			role_id
		FROM users WHERE id = $1;
	`
	user := new(users.User)
	if err := u.db.Get(user, query, userId); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	return user, nil
}

func (u *usersrepository) DeleteOauth(oauthId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	query := `
		DELETE FROM oauth WHERE id =$1;
	`
	if _, err := u.db.ExecContext(ctx, query, oauthId); err != nil {
		return fmt.Errorf("oauth id not found")
	}

	return nil
}
