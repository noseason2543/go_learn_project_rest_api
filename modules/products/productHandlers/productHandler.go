package productHandlers

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/appInfo"
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/files"
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
	insertProductErr  productsHandlersErrCode = "products-003"
	updateProductErr  productsHandlersErrCode = "products-004"
	deleteProductErr  productsHandlersErrCode = "products-005"
)

type IProductHandler interface {
	FindOneProduct(fiber.Ctx) error
	FindProduct(fiber.Ctx) error
	AddProduct(fiber.Ctx) error
	UpdateProduct(fiber.Ctx) error
	DeleteProduct(fiber.Ctx) error
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

func (h *productHandler) AddProduct(c fiber.Ctx) error {
	req := &products.Product{
		Category: &appInfo.Category{},
		Images:   make([]*entities.Image, 0),
	}
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(insertProductErr),
			err.Error(),
		).Res()
	}
	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(insertProductErr),
			"category id is invalid",
		).Res()
	}

	product, err := h.productUsecase.AddProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).SuccessResponse(fiber.StatusCreated, product).Res()
}

func (h *productHandler) UpdateProduct(c fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")
	req := &products.Product{
		Images:   make([]*entities.Image, 0),
		Category: &appInfo.Category{},
	}
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateProductErr),
			err.Error(),
		).Res()
	}
	req.Id = productId

	product, err := h.productUsecase.UpdateProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).SuccessResponse(fiber.StatusOK, product).Res()
}

func (h *productHandler) DeleteProduct(c fiber.Ctx) error {
	req := strings.Trim(c.Params("product_id"), "")

	product, err := h.productUsecase.FindOneProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, image := range product.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: "images/test/",
			FileName:    image.FileName,
		})
	}

	if err := h.fileUsecases.DeleteFileOnGCP(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	if err := h.productUsecase.DeleteProduct(product.Id); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).SuccessResponse(fiber.StatusNoContent, nil).Res()
}
