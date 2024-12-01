package usersHandlers

import (
	"fmt"
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/users"
	"go_learn_project_rest_api/modules/users/usersUsecases"
	"go_learn_project_rest_api/pkgs/auth"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type userHandlersErrCode string

const (
	signUpCustomerErrCode     userHandlersErrCode = "users-001"
	signInErrCode             userHandlersErrCode = "users-002"
	refreshPassportErrCode    userHandlersErrCode = "users-003"
	signOutErrCode            userHandlersErrCode = "users-004"
	signUpAdminErrCode        userHandlersErrCode = "users-005"
	generateAdminTokenErrCode userHandlersErrCode = "users-006"
	GetUserProfileErrCode     userHandlersErrCode = "users-007"
)

type IUsersHandlers interface {
	SignUpCustomer(fiber.Ctx) error
	SignIn(fiber.Ctx) error
	RefreshPassport(fiber.Ctx) error
	SignOut(fiber.Ctx) error
	SignUpAdmin(fiber.Ctx) error
	GenerateAdminToken(fiber.Ctx) error
	GetUserProfile(fiber.Ctx) error
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

func (h *usersHandlers) SignUpAdmin(c fiber.Ctx) error {
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

	result, err := h.userUsecases.InsertAdmin(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signUpCustomerErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).SuccessResponse(fiber.StatusCreated, result).Res()
}

func (h *usersHandlers) SignIn(c fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.Bind().Body(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signInErrCode),
			err.Error(),
		).Res()
	}

	passport, err := h.userUsecases.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signInErrCode),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).SuccessResponse(fiber.StatusOK, passport).Res()
}

func (h *usersHandlers) RefreshPassport(c fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	fmt.Println("qwe: ", string(c.Body()))
	if err := c.Bind().Body(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(refreshPassportErrCode),
			err.Error(),
		).Res()
	}

	passport, err := h.userUsecases.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(refreshPassportErrCode),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).SuccessResponse(fiber.StatusOK, passport).Res()
}

func (h *usersHandlers) SignOut(c fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signOutErrCode),
			err.Error(),
		).Res()
	}

	if err := h.userUsecases.DeleteOauth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signOutErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).SuccessResponse(fiber.StatusOK, nil).Res()
}

func (h *usersHandlers) GenerateAdminToken(c fiber.Ctx) error {
	adminToken, err := auth.NewAuth(auth.Admin, h.cfg.Jwt(), nil)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(generateAdminTokenErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).SuccessResponse(
		fiber.StatusOK,
		&struct {
			Token string `json:"token"`
		}{
			Token: adminToken.SignToken(),
		},
	).Res()
}

func (h *usersHandlers) GetUserProfile(c fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")
	result, err := h.userUsecases.GetUserProfile(userId)
	if err != nil {
		if strings.HasPrefix(err.Error(), "get user failed: sql: no rows in result set") {
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(GetUserProfileErrCode),
				err.Error(),
			).Res()
		} else {
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(GetUserProfileErrCode),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).SuccessResponse(fiber.StatusOK, result).Res()
}
