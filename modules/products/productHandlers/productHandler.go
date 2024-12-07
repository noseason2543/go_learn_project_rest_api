package productHandlers

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/files/fileUsecases"
	"go_learn_project_rest_api/modules/products/productUsecases"
)

type IProductHandler interface {
}

type productHandler struct {
	productUsecase productUsecases.IProductUsecases
	cfg            config.IConfig
	fileUsecases   fileUsecases.IFileUsecases
}

func ProductHandler(productUsecase productUsecases.IProductUsecases, cfg config.IConfig, fileUsecases fileUsecases.IFileUsecases) IProductHandler {
	return &productHandler{
		productUsecase: productUsecase,
		cfg:            cfg,
		fileUsecases:   fileUsecases,
	}
}
