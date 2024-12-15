package orderHandlers

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/orders"
	"go_learn_project_rest_api/modules/orders/orderUsecases"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

type ordersHandlersErrCode string

const (
	findOneOrderErr ordersHandlersErrCode = "orders-001"
	findOrderErr    ordersHandlersErrCode = "orders-002"
	insertOrderErr  ordersHandlersErrCode = "orders-003"
	updateOrderErr  ordersHandlersErrCode = "orders-004"
)

type IOrderHandlers interface {
	FindOneOrder(fiber.Ctx) error
	FindOrder(fiber.Ctx) error
}

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

func (h *orderHandlers) FindOneOrder(c fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")

	order, err := h.orderUsecases.FindOneOrder(orderId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).SuccessResponse(fiber.StatusOK, order).Res()
}

func (h *orderHandlers) FindOrder(c fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
	}
	if err := c.Bind().Query(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findOrderErr),
			err.Error(),
		).Res()
	}

	// Paginate
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	// Sort
	orderByMap := map[string]string{
		"id":         `"o"."id"`,
		"created_at": `"o"."created_at"`,
	}
	if orderByMap[req.OrderBy] == "" {
		req.OrderBy = orderByMap["id"]
	}

	req.Sort = strings.ToUpper(req.Sort)
	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}
	if sortMap[req.Sort] == "" {
		req.Sort = sortMap["DESC"]
	}

	// Date	YYYY-MM-DD
	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(findOrderErr),
				"start date is invalid",
			).Res()
		}
		req.StartDate = start.Format("2006-01-02")
	}
	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(findOrderErr),
				"end date is invalid",
			).Res()
		}
		req.EndDate = end.Format("2006-01-02")
	}

	return entities.NewResponse(c).SuccessResponse(
		fiber.StatusOK,
		h.orderUsecases.FindOrder(req),
	).Res()
}
