package servers

import (
	"encoding/json"
	"go_learn_project_rest_api/config"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

type IServer interface {
	Start()
	// GetServer() *server
}

type server struct {
	app *fiber.App
	db  *sqlx.DB
	cfg config.IConfig
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		cfg: cfg,
		db:  db,
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

func (s *server) Start() {
	//middlewares
	middlewares := InitMiddlewares(s)
	s.app.Use(middlewares.Logger())
	s.app.Use(middlewares.Cors())
	//module
	v1 := s.app.Group("v1")
	modules := InitModule(v1, s, middlewares)
	modules.MonitorModule()

	s.app.Use(middlewares.RouterCheck())
	//graceful shut down
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("server is shutting down . . .")
		_ = s.app.Shutdown()
	}()

	// listen to host:port
	log.Printf("server start on %v", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())
}
