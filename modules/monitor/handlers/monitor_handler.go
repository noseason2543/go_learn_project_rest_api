package handlers

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/monitor"

	"github.com/gofiber/fiber/v3"
)

// handler only have context parameter and error return type
type IMonitorHandler interface {
	HealthCheck(c fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.IConfig
}

func MonitorHandler(cfg config.IConfig) IMonitorHandler {
	return &monitorHandler{
		cfg: cfg,
	}
}

func (h *monitorHandler) HealthCheck(c fiber.Ctx) error {
	res := &monitor.Monitor{
		Name:    h.cfg.App().Name(),
		Version: h.cfg.App().Version(),
	}

	return entities.NewResponse(c).SuccessResponse(fiber.StatusOK, res).Res()
}
