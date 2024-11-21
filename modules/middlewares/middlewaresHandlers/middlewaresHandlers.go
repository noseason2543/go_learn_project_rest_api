package middlewaresHandler

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/middlewares/middlewaresUsecases"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type middlewaresHandlerErrCode string

const (
	routerCheckErr middlewaresHandlerErrCode = "middlewares-001"
)

type IMiddlewaresHandlers interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
}

type middlewaresHandlers struct {
	cfg                config.IConfig
	middlewareUsecases middlewaresUsecases.IMiddlewaresUsecases
}

func MiddlewaresHandlers(cfg config.IConfig, m middlewaresUsecases.IMiddlewaresUsecases) IMiddlewaresHandlers {
	return &middlewaresHandlers{
		cfg:                cfg,
		middlewareUsecases: m,
	}
}

func (h *middlewaresHandlers) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "PATCH", "DELETE"},
		AllowHeaders:     []string{""},
		AllowCredentials: false,
		ExposeHeaders:    []string{""},
		MaxAge:           0,
	})
}

func (h *middlewaresHandlers) RouterCheck() fiber.Handler {
	return func(c fiber.Ctx) error {
		return entities.NewResponse(c).Error(fiber.StatusNotFound, string(routerCheckErr), "router not found").Res()
	}
}

func (h *middlewaresHandlers) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}
