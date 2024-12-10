package productRepositories

import (
	"context"
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
	InsertProduct(*products.Product) (*products.Product, error)
	UpdateProduct(*products.Product) (*products.Product, error)
	DeleteProduct(string) error
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

func (r *productRepository) InsertProduct(req *products.Product) (*products.Product, error) {
	builder := productPatterns.InsertProductBuilder(r.db, req)
	productId, err := productPatterns.InsertProductEngineer(builder).InsertProduct()
	if err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) UpdateProduct(req *products.Product) (*products.Product, error) {
	builder := productPatterns.UpdateProductBuilder(r.db, req, r.fileUsecases)
	engineer := productPatterns.UpdateProductEngineer(builder)

	if err := engineer.UpdateProduct(); err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(req.Id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) DeleteProduct(productId string) error {
	ctx := context.Background()
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	query := `DELETE FROM products WHERE id = $1;`

	if _, err := tx.ExecContext(ctx, query, productId); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete product failed: %v", err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
