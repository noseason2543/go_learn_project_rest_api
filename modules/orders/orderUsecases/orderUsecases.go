package orderUsecases

import (
	"go_learn_project_rest_api/modules/orders/orderRepositories"
	"go_learn_project_rest_api/modules/products/productRepositories"
)

type IOrderUsecases interface{}

type orderUsecases struct {
	orderRepository   orderRepositories.IOrderRepository
	productRepository productRepositories.IProductRepository
}

func OrderUsecases(orderRepository orderRepositories.IOrderRepository, productRepository productRepositories.IProductRepository) IOrderUsecases {
	return &orderUsecases{
		orderRepository:   orderRepository,
		productRepository: productRepository,
	}
}
