package productHandlers

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/files/fileUsecases"
	"go_learn_project_rest_api/modules/products"
	"go_learn_project_rest_api/modules/products/productUsecases"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type productsHandlersErrCode string

const (
	findOneProductErr productsHandlersErrCode = "products-001"
	findProductErr    productsHandlersErrCode = "products-002"
)

type IProductHandler interface {
	FindOneProduct(fiber.Ctx) error
	FindProduct(fiber.Ctx) error
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

func (h *productHandler) FindOneProduct(c fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(findOneProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).SuccessResponse(fiber.StatusOK, product).Res()
}

func (h *productHandler) FindProduct(c fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.Bind().Query(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findProductErr),
			err.Error(),
		).Res()
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}
	if req.Sort == "" {
		req.Sort = "ASC"
	}

	products := h.productUsecase.FindProduct(req)
	return entities.NewResponse(c).SuccessResponse(fiber.StatusOK, products).Res()
}
