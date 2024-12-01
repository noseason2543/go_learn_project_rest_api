package usersPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"go_learn_project_rest_api/modules/users"
	"time"

	"github.com/jmoiron/sqlx"
)

type IInsertUser interface {
	Customer() (IInsertUser, error)
	Admin() (IInsertUser, error)
	Result() (*users.UserPassport, error)
}

type userReq struct {
	id  string
	req *users.UserRegisterReq
	db  *sqlx.DB
}

type customer struct {
	*userReq
}

type admin struct {
	*userReq
}

func InsertUser(db *sqlx.DB, req *users.UserRegisterReq, isAdmin bool) IInsertUser {
	if isAdmin {
		return newAdmin(db, req)
	}
	return newCustomer(db, req)
}

func newCustomer(db *sqlx.DB, req *users.UserRegisterReq) IInsertUser {
	return &customer{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func newAdmin(db *sqlx.DB, req *users.UserRegisterReq) IInsertUser {
	return &admin{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func (f *userReq) Customer() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		INSERT INTO users (
			email,
			username,
			password,
			role_id
		) VALUES (
			$1, $2, $3, 1
		)
		RETURNING id
	`
	if err := f.db.QueryRowContext(
		ctx,
		query,
		f.req.Email,
		f.req.Username,
		f.req.Password,
	).Scan(&f.id); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return f, nil
}

func (f *userReq) Admin() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	query := `
		INSERT INTO users (
			email,
			username,
			password,
			role_id
		) VALUES (
			$1, $2, $3, 2
		)
		RETURNING id
	`
	if err := f.db.QueryRowContext(
		ctx,
		query,
		f.req.Email,
		f.req.Username,
		f.req.Password,
	).Scan(&f.id); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return f, nil
}

func (f *userReq) Result() (*users.UserPassport, error) {
	query := `
		SELECT json_build_object(
			'user', json_build_object(
				'id', id,
				'email', email,
				'username', username,
				'role_id', role_id
			),
			'token', NULL
		)
		FROM users 
		WHERE id = $1
	`

	var jsonData []byte
	if err := f.db.Get(&jsonData, query, f.id); err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %w", err)
	}

	user := new(users.UserPassport)
	if err := json.Unmarshal(jsonData, user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}

	return user, nil
}
