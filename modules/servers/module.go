package servers

import (
	"go_learn_project_rest_api/modules/appInfo/appInfoHandlers"
	"go_learn_project_rest_api/modules/appInfo/appInfoRepositories"
	"go_learn_project_rest_api/modules/appInfo/appInfoUsecases"
	middlewaresHandler "go_learn_project_rest_api/modules/middlewares/middlewaresHandlers"
	"go_learn_project_rest_api/modules/middlewares/middlewaresRepository"
	"go_learn_project_rest_api/modules/middlewares/middlewaresUsecases"
	"go_learn_project_rest_api/modules/monitor/handlers"
	"go_learn_project_rest_api/modules/users/usersHandlers"
	"go_learn_project_rest_api/modules/users/usersRepositories"
	"go_learn_project_rest_api/modules/users/usersUsecases"

	"github.com/gofiber/fiber/v3"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppInfoModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
	mid    middlewaresHandler.IMiddlewaresHandlers
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandler.IMiddlewaresHandlers) IModuleFactory {
	return &moduleFactory{
		router: r,
		server: s,
		mid:    mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandler.IMiddlewaresHandlers {
	repository := middlewaresRepository.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecases(repository)
	handler := middlewaresHandler.MiddlewaresHandlers(s.cfg, usecase)
	return handler

}

func (m *moduleFactory) MonitorModule() {
	handler := handlers.MonitorHandler(m.server.cfg)
	m.router.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.server.db)
	usecase := usersUsecases.UsersUsecases(m.server.cfg, repository)
	handlers := usersHandlers.UsersHandlers(m.server.cfg, usecase)

	router := m.router.Group("/users")
	router.Post("/signup", handlers.SignUpCustomer, m.mid.ApiKeyAuth())
	router.Post("/signin", handlers.SignIn, m.mid.ApiKeyAuth())
	router.Post("/refresh", handlers.RefreshPassport, m.mid.ApiKeyAuth())
	router.Post("/signout", handlers.SignOut, m.mid.ApiKeyAuth())
	router.Post("/signup-admin", handlers.SignUpAdmin, m.mid.JwtAuth(), m.mid.Authorize(2))

	router.Get("/admin/secret", handlers.GenerateAdminToken, m.mid.JwtAuth(), m.mid.Authorize(2))
	router.Get("/profile/:user_id", handlers.GetUserProfile, m.mid.JwtAuth(), m.mid.ParamsCheck())

}

func (m *moduleFactory) AppInfoModule() {
	repository := appInfoRepositories.AppInfoRepository(m.server.db)
	usecase := appInfoUsecases.AppInfoUsecases(repository)
	handlers := appInfoHandlers.AppInfoHandler(m.server.cfg, usecase)
	_ = handlers

	router := m.router.Group("/appinfo")

	router.Get("/apikey", handlers.GenerateApiKey, m.mid.JwtAuth(), m.mid.Authorize(2))
	router.Get("/category", handlers.FindCategory, m.mid.ApiKeyAuth())
}
