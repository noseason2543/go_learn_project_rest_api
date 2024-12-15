package orderUsecases

import (
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/orders"
	"go_learn_project_rest_api/modules/orders/orderRepositories"
	"go_learn_project_rest_api/modules/products/productRepositories"
	"math"
)

type IOrderUsecases interface {
	FindOneOrder(string) (*orders.Order, error)
	FindOrder(*orders.OrderFilter) *entities.PaginateRes
}

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

func (u *orderUsecases) FindOneOrder(orderId string) (*orders.Order, error) {
	return u.orderRepository.FindOneOrder(orderId)
}

func (u *orderUsecases) FindOrder(req *orders.OrderFilter) *entities.PaginateRes {
	orders, count := u.orderRepository.FindOrder(req)
	return &entities.PaginateRes{
		Data:      orders,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}
