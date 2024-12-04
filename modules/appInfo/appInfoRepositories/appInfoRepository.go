package appInfoRepositories

import (
	"fmt"
	"go_learn_project_rest_api/modules/appInfo"
	"strings"

	"github.com/jmoiron/sqlx"
)

type IAppInfoRepository interface {
	FindCategory(*appInfo.CategoryFilter) ([]*appInfo.Category, error)
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
