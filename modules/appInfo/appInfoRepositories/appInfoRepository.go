package appInfoRepositories

import (
	"context"
	"fmt"
	"go_learn_project_rest_api/modules/appInfo"
	"strings"

	"github.com/jmoiron/sqlx"
)

type IAppInfoRepository interface {
	FindCategory(*appInfo.CategoryFilter) ([]*appInfo.Category, error)
	InsertCategory([]*appInfo.Category) error
	DeleteCategory(int) error
}

type appInfoRepository struct {
	db *sqlx.DB
}

func AppInfoRepository(db *sqlx.DB) IAppInfoRepository {
	return &appInfoRepository{
		db: db,
	}
}

func (r *appInfoRepository) FindCategory(req *appInfo.CategoryFilter) ([]*appInfo.Category, error) {
	query := `
        SELECT 
            id,
            title
        FROM categories
    `

	filterValue := make([]any, 0)
	if req.Title != "" {
		query += `
            WHERE LOWER(title) LIKE $1
        `
		filterValue = append(filterValue, "%"+strings.ToLower(req.Title)+"%")
	}

	category := make([]*appInfo.Category, 0)
	if err := r.db.Select(&category, query, filterValue...); err != nil {
		return nil, fmt.Errorf("select category failed: %v", err)
	}

	return category, nil
}

func (r *appInfoRepository) InsertCategory(req []*appInfo.Category) error {
	ctx := context.Background()
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO categories (
            title
        ) VALUES
    `

	valuesStack := make([]any, 0)
	for i, cate := range req {
		valuesStack = append(valuesStack, cate.Title)

		if i != len(req)-1 {
			query += fmt.Sprintf(`
		($%d),`, i+1)
		} else {
			query += fmt.Sprintf(`
		($%d)`, i+1)
		}
	}
	query += `
        RETURNING id;
    `

	rows, err := tx.QueryxContext(ctx, query, valuesStack...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("insert category failed: %v", err)
	}

	var i int
	for rows.Next() {
		if err := rows.Scan(&req[i].Id); err != nil {
			return fmt.Errorf("scan categories id failed: %v", err)
		}
		i++
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *appInfoRepository) DeleteCategory(id int) error {
	ctx := context.Background()
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	query := `DELETE FROM categories WHERE id = $1;`

	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete categories failed: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
