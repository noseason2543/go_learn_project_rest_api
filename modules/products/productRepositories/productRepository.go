package productRepositories

import (
	"encoding/json"
	"fmt"
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/files/fileUsecases"
	"go_learn_project_rest_api/modules/products"
	"go_learn_project_rest_api/modules/products/productPatterns"

	"github.com/jmoiron/sqlx"
)

type IProductRepository interface {
	FindOneProduct(string) (*products.Product, error)
	FindProduct(*products.ProductFilter) ([]*products.Product, int)
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
                p.price,
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
                p.updated_at,
                (
                    SELECT 
                        COALESCE(array_to_json(array_agg(it)), '[]'::json)
                    FROM (
                        SELECT 
                            i.id,
                            i.filename,
                            i.url
                        FROM images i
                        WHERE i.product_id = p.id
                    ) AS it
                ) AS images
            FROM products p
            WHERE p.id = $1
            LIMIT 1
        ) AS t;
    `
	bytesResult := make([]byte, 0)
	result := &products.Product{
		Images: make([]*entities.Image, 0),
	}

	if err := r.db.Get(&bytesResult, query, productId); err != nil {
		return nil, fmt.Errorf("get product failed: %v", err)
	}

	if err := json.Unmarshal(bytesResult, &result); err != nil {
		return nil, fmt.Errorf("unmarshal product failed: %v", err)
	}

	return result, nil
}

func (r *productRepository) FindProduct(req *products.ProductFilter) ([]*products.Product, int) {
	builder := productPatterns.FindProductBuilder(r.db, req)
	engineer := productPatterns.FindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()
	return result, count
}
