package productRepositories

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/files/fileUsecases"
	"go_learn_project_rest_api/modules/products"

	"github.com/jmoiron/sqlx"
)

type IProductRepository interface {
}

type productRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	fileUsecases fileUsecases.IFileUsecases
}

func ProductRepository(db *sqlx.DB, cfg config.IConfig, fileUsecases fileUsecases.IFileUsecases) IProductRepository {
	return &productRepository{
		db:           db,
		cfg:          cfg,
		fileUsecases: fileUsecases,
	}
}

func (r *productRepository) FindOneProduct(productId string) (*products.Product, error) {
	query := `
        SELECT
            to_jsonb(t)
        FROM (
            SELECT
                p.id,
                p.title,
                p.description,
                (
                    SELECT
                        to_jsonb(ct)
                    FROM (
                        SELECT 
                            c.id,
                            c.title
                        FROM categories c LEFT JOIN products_categories pc ON pc.category_id = c.id
                        WHERE pc.product_id = p.id
                    ) AS ct
                ) AS categories,
                p.created_at,
                p.updated_at
            FROM products p
        ) AS t
    `

	_ = query
	return nil, nil
}
