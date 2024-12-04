package appInfoHandlers

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/appInfo"
	"go_learn_project_rest_api/modules/appInfo/appInfoUsecases"
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/pkgs/auth"

	"github.com/gofiber/fiber/v3"
)

type appInfoHandlerErrCode string

const (
	generateApiKeyTokenErrCode appInfoHandlerErrCode = "appInfo-001"
	findCategoryErrCode        appInfoHandlerErrCode = "appInfo-002"
)

type IAppInfoHandler interface {
	GenerateApiKey(fiber.Ctx) error
	FindCategory(fiber.Ctx) error
}

type appInfoHandler struct {
	cfg             config.IConfig
	appInfoUsecases appInfoUsecases.IAppInfoUsecases
}

func AppInfoHandler(cfg config.IConfig, usecase appInfoUsecases.IAppInfoUsecases) IAppInfoHandler {
	return &appInfoHandler{
		cfg:             cfg,
		appInfoUsecases: usecase,
	}
}

func (h *appInfoHandler) GenerateApiKey(c fiber.Ctx) error {
	apiKey, err := auth.NewAuth(auth.ApiKey, h.cfg.Jwt(), nil)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(generateApiKeyTokenErrCode),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).SuccessResponse(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()
}

func (h *appInfoHandler) FindCategory(c fiber.Ctx) error {
	req := new(appInfo.CategoryFilter)
	if err := c.Bind().Query(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(findCategoryErrCode),
			err.Error(),
		).Res()
	}

	category, err := h.appInfoUsecases.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(findCategoryErrCode),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).SuccessResponse(fiber.StatusOK, category).Res()
}
