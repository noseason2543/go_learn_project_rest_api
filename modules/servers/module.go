package servers

import (
	"go_learn_project_rest_api/modules/monitor/handlers"

	"github.com/gofiber/fiber/v3"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
}

func InitModule(r fiber.Router, s *server) IModuleFactory {
	return &moduleFactory{
		router: r,
		server: s,
	}
}

func (m *moduleFactory) MonitorModule() {
	handler := handlers.MonitorHandler(m.server.cfg)
	m.router.Get("/", handler.HealthCheck)
}
