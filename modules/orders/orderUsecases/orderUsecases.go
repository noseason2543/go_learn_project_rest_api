package orderUsecases

import (
	"fmt"
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/orders"
	"go_learn_project_rest_api/modules/orders/orderRepositories"
	"go_learn_project_rest_api/modules/products/productRepositories"
	"go_learn_project_rest_api/pkgs/utils"
	"math"
)

type IOrderUsecases interface {
	FindOneOrder(string) (*orders.Order, error)
	FindOrder(*orders.OrderFilter) *entities.PaginateRes
	InsertOrder(*orders.Order) (*orders.Order, error)
	UpdateOrder(*orders.Order) (*orders.Order, error)
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

func (u *orderUsecases) InsertOrder(req *orders.Order) (*orders.Order, error) {
	// Check if products is exists
	for i := range req.Products {
		if req.Products[i].Product == nil {
			return nil, fmt.Errorf("product is nil")
		}

		prod, err := u.productRepository.FindOneProduct(req.Products[i].Product.Id)
		if err != nil {
			return nil, err
		}
		utils.Debug(prod)

		// Set price
		req.TotalPaid += req.Products[i].Product.Price * float64(req.Products[i].Qty)
		req.Products[i].Product = prod
	}

	orderId, err := u.orderRepository.InsertOrder(req)
	if err != nil {
		return nil, err
	}

	order, err := u.orderRepository.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *orderUsecases) UpdateOrder(req *orders.Order) (*orders.Order, error) {
	if err := u.orderRepository.UpdateOrder(req); err != nil {
		return nil, err
	}

	order, err := u.orderRepository.FindOneOrder(req.Id)
	if err != nil {
		return nil, err
	}
	return order, nil
}
