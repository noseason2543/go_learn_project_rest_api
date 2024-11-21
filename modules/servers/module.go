package servers

import (
	middlewaresHandler "go_learn_project_rest_api/modules/middlewares/middlewaresHandlers"
	"go_learn_project_rest_api/modules/middlewares/middlewaresRepository"
	"go_learn_project_rest_api/modules/middlewares/middlewaresUsecases"
	"go_learn_project_rest_api/modules/monitor/handlers"

	"github.com/gofiber/fiber/v3"
)

type IModuleFactory interface {
	MonitorModule()
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
