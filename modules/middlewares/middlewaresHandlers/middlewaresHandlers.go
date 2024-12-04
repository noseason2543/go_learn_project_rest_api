package middlewaresHandler

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/middlewares/middlewaresUsecases"
	"go_learn_project_rest_api/pkgs/auth"
	"go_learn_project_rest_api/pkgs/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type middlewaresHandlerErrCode string

const (
	routerCheckErr middlewaresHandlerErrCode = "middlewares-001"
	jwtAuthErr     middlewaresHandlerErrCode = "middlewares-002"
	paramsCheckErr middlewaresHandlerErrCode = "middlewares-003"
	authorizeErr   middlewaresHandlerErrCode = "middlewares-004"
	apiKeyErr      middlewaresHandlerErrCode = "middlewares-005"
)

type IMiddlewaresHandlers interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(...int) fiber.Handler
	ApiKeyAuth() fiber.Handler
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

func (h *middlewaresHandlers) JwtAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("authorization"), "Bearer ")
		result, err := auth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.StatusUnauthorized,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims

		if !h.middlewareUsecases.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.StatusUnauthorized,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}
		c.Locals("userId", claims.Id)
		c.Locals("roleId", claims.RoleId)
		return c.Next()
	}
}

func (h *middlewaresHandlers) ParamsCheck() fiber.Handler {
	return func(c fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.StatusUnauthorized,
				string(paramsCheckErr),
				"userId not match",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewaresHandlers) Authorize(expectedRole ...int) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRoleId, ok := c.Locals("roleId").(int)
		if !ok {
			return entities.NewResponse(c).Error(
				fiber.StatusUnauthorized,
				string(authorizeErr),
				"role_id is not int type",
			).Res()
		}

		roles, err := h.middlewareUsecases.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(authorizeErr),
				err.Error(),
			).Res()
		}

		sum := 0
		for _, val := range expectedRole {
			sum += val
		}
		expectedValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))
		for i, val := range userValueBinary {
			if val&expectedValueBinary[i] == 1 {
				return c.Next()
			}
		}
		return entities.NewResponse(c).Error(
			fiber.StatusUnauthorized,
			string(authorizeErr),
			"no permission to access",
		).Res()
	}
}

func (h *middlewaresHandlers) ApiKeyAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Get("X-Api-Key")
		_, err := auth.ParseApiKeyToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.StatusUnauthorized,
				string(apiKeyErr),
				"apiKey is invalided",
			).Res()
		}

		return c.Next()
	}
}
