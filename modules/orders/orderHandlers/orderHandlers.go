package orderHandlers

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/orders/orderUsecases"
)

type IOrderHandlers interface{}

type orderHandlers struct {
	cfg           config.IConfig
	orderUsecases orderUsecases.IOrderUsecases
}

func OrderHandlers(cfg config.IConfig, orderUsecase orderUsecases.IOrderUsecases) IOrderHandlers {
	return &orderHandlers{
		cfg:           cfg,
		orderUsecases: orderUsecase,
	}
}
