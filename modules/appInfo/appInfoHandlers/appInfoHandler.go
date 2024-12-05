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
	insertCategoryErrCode      appInfoHandlerErrCode = "appInfo-003"
	deleteCategoryErrCode      appInfoHandlerErrCode = "appInfo-004"
)

type IAppInfoHandler interface {
	GenerateApiKey(fiber.Ctx) error
	FindCategory(fiber.Ctx) error
	InsertCategory(fiber.Ctx) error
	DeleteCategory(fiber.Ctx) error
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

func (h *appInfoHandler) InsertCategory(c fiber.Ctx) error {
	req := make([]*appInfo.Category, 0)
	if err := c.Bind().Body(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(insertCategoryErrCode),
			err.Error(),
		).Res()
	}

	if len(req) == 0 {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(insertCategoryErrCode),
			"categories request are empty",
		).Res()
	}

	if err := h.appInfoUsecases.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(insertCategoryErrCode),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).SuccessResponse(fiber.StatusCreated, req).Res()
}

func (h *appInfoHandler) DeleteCategory(c fiber.Ctx) error {
	req := new(appInfo.RequestCategoryId)
	if err := c.Bind().Body(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(deleteCategoryErrCode),
			err.Error(),
		).Res()
	}

	if err := h.appInfoUsecases.DeleteCategory(req.Id); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(deleteCategoryErrCode),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).SuccessResponse(fiber.StatusNoContent, nil).Res()
}
