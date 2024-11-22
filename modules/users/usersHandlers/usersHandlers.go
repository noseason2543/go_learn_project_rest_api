package usersHandlers

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/users"
	"go_learn_project_rest_api/modules/users/usersUsecases"

	"github.com/gofiber/fiber/v3"
)

type userHandlersErrCode string

const (
	signUpCustomerErrCode userHandlersErrCode = "users-001"
)

type IUsersHandlers interface {
	SignUpCustomer(fiber.Ctx) error
}

type usersHandlers struct {
	cfg          config.IConfig
	userUsecases usersUsecases.IUsersUsecases
}

func UsersHandlers(cfg config.IConfig, userUsecases usersUsecases.IUsersUsecases) IUsersHandlers {
	return &usersHandlers{
		cfg:          cfg,
		userUsecases: userUsecases,
	}
}

func (h *usersHandlers) SignUpCustomer(c fiber.Ctx) error {
	req := new(users.UserRegisterReq)
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signUpCustomerErrCode),
			err.Error(),
		).Res()
	}

	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signUpCustomerErrCode),
			"email patterns is invalid",
		).Res()
	}

	result, err := h.userUsecases.InsertCustomer(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signUpCustomerErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).SuccessResponse(fiber.StatusCreated, result).Res()
}
