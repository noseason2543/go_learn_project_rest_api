package productUsecases

import "go_learn_project_rest_api/modules/products/productRepositories"

type IProductUsecases interface {
}

type productUsecases struct {
	productRepositories productRepositories.IProductRepository
}

func ProductUsecases(productRepositories productRepositories.IProductRepository) IProductUsecases {
	return &productUsecases{
		productRepositories: productRepositories,
	}
}
