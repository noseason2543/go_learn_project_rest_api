package productUsecases

import (
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/products"
	"go_learn_project_rest_api/modules/products/productRepositories"
	"math"
)

type IProductUsecases interface {
	FindOneProduct(string) (*products.Product, error)
	FindProduct(*products.ProductFilter) *entities.PaginateRes
	AddProduct(*products.Product) (*products.Product, error)
	UpdateProduct(*products.Product) (*products.Product, error)
	DeleteProduct(string) error
}

type productUsecases struct {
	productRepositories productRepositories.IProductRepository
}

func ProductUsecases(productRepositories productRepositories.IProductRepository) IProductUsecases {
	return &productUsecases{
		productRepositories: productRepositories,
	}
}

func (u *productUsecases) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.productRepositories.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productUsecases) FindProduct(req *products.ProductFilter) *entities.PaginateRes {
	products, count := u.productRepositories.FindProduct(req)
	return &entities.PaginateRes{
		Data:      products,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
		TotalItem: count,
	}

}

func (u *productUsecases) AddProduct(req *products.Product) (*products.Product, error) {
	return u.productRepositories.InsertProduct(req)
}

func (u *productUsecases) UpdateProduct(req *products.Product) (*products.Product, error) {
	return u.productRepositories.UpdateProduct(req)
}

func (u *productUsecases) DeleteProduct(productId string) error {
	return u.productRepositories.DeleteProduct(productId)
}
