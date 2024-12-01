package middlewaresRepository

import (
	"fmt"
	"go_learn_project_rest_api/modules/middlewares"

	"github.com/jmoiron/sqlx"
)

type IMiddlewaresRepository interface {
	FindAccessToken(userId, token string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresRepository struct {
	db *sqlx.DB
}

func MiddlewaresRepository(db *sqlx.DB) IMiddlewaresRepository {
	return &middlewaresRepository{
		db: db,
	}
}

func (r *middlewaresRepository) FindAccessToken(userId, token string) bool {
	query := `
		SELECT (CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
		FROM oauth WHERE user_id = $1 AND access_token = $2;
	`
	var check bool
	if err := r.db.Get(&check, query, userId, token); err != nil {
		return false
	}
	return check
}

func (r *middlewaresRepository) FindRole() ([]*middlewares.Role, error) {
	query := `
		SELECT id, title FROM roles ORDER BY id DESC
	`
	roles := make([]*middlewares.Role, 0)
	if err := r.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("roles are empty")
	}
	return roles, nil
}
