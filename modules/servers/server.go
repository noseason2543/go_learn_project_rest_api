package servers

import (
	"encoding/json"
	"go_learn_project_rest_api/config"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

type IServer interface {
	Start()
	// GetServer() *server
}

type server struct {
	app    *fiber.App
	db     *sqlx.DB
	config config.IConfig
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		config: cfg,
		db:     db,
		app: fiber.New(
			fiber.Config{
				AppName:      cfg.App().Name(),
				BodyLimit:    cfg.App().BodyLimit(),
				ReadTimeout:  cfg.App().ReadTimeout(),
				WriteTimeout: cfg.App().WriteTimeout(),
				JSONEncoder:  json.Marshal, // make fiber faster
				JSONDecoder:  json.Unmarshal,
			},
		),
	}
}

func (s *server) Start() {}
