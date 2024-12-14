package servers

import (
	"go_learn_project_rest_api/modules/appInfo/appInfoHandlers"
	"go_learn_project_rest_api/modules/appInfo/appInfoRepositories"
	"go_learn_project_rest_api/modules/appInfo/appInfoUsecases"
	"go_learn_project_rest_api/modules/files/fileHandlers"
	"go_learn_project_rest_api/modules/files/fileUsecases"
	middlewaresHandler "go_learn_project_rest_api/modules/middlewares/middlewaresHandlers"
	"go_learn_project_rest_api/modules/middlewares/middlewaresRepository"
	"go_learn_project_rest_api/modules/middlewares/middlewaresUsecases"
	"go_learn_project_rest_api/modules/monitor/handlers"
	"go_learn_project_rest_api/modules/orders/orderHandlers"
	"go_learn_project_rest_api/modules/orders/orderRepositories"
	"go_learn_project_rest_api/modules/orders/orderUsecases"
	"go_learn_project_rest_api/modules/products/productHandlers"
	"go_learn_project_rest_api/modules/products/productRepositories"
	"go_learn_project_rest_api/modules/products/productUsecases"
	"go_learn_project_rest_api/modules/users/usersHandlers"
	"go_learn_project_rest_api/modules/users/usersRepositories"
	"go_learn_project_rest_api/modules/users/usersUsecases"

	"github.com/gofiber/fiber/v3"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppInfoModule()
	FilesModule()
	ProductModule()
	OrderModule()
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

	router := m.router.Group("/appinfo")

	router.Get("/apikey", handlers.GenerateApiKey, m.mid.JwtAuth(), m.mid.Authorize(2))
	router.Post("/insertcategory", handlers.InsertCategory, m.mid.JwtAuth(), m.mid.Authorize(2))
	router.Post("/deletecategory", handlers.DeleteCategory, m.mid.JwtAuth(), m.mid.Authorize(2))
	router.Get("/category", handlers.FindCategory, m.mid.ApiKeyAuth())
}

func (m *moduleFactory) FilesModule() {
	usecase := fileUsecases.FileUsecases(m.server.cfg)
	handlers := fileHandlers.FileHandlers(m.server.cfg, usecase)

	router := m.router.Group("/files")
	router.Post("/upload", handlers.UploadFiles, m.mid.JwtAuth(), m.mid.Authorize(2))
	router.Post("/delete", handlers.DeleteFile, m.mid.JwtAuth(), m.mid.Authorize(2))
}

func (m *moduleFactory) ProductModule() {
	fileUsecase := fileUsecases.FileUsecases(m.server.cfg)
	repository := productRepositories.ProductRepository(m.server.db, m.server.cfg, fileUsecase)
	usecase := productUsecases.ProductUsecases(repository)
	handlers := productHandlers.ProductHandler(usecase, m.server.cfg, fileUsecase)

	router := m.router.Group("/products")
	router.Post("/addProduct", handlers.AddProduct, m.mid.JwtAuth(), m.mid.Authorize(2))
	router.Patch("/:product_id", handlers.UpdateProduct, m.mid.JwtAuth(), m.mid.Authorize(2))
	router.Get("/", handlers.FindProduct, m.mid.ApiKeyAuth())
	router.Get("/:product_id", handlers.FindOneProduct, m.mid.ApiKeyAuth())

	router.Delete("/:product_id", handlers.DeleteProduct, m.mid.JwtAuth(), m.mid.Authorize(2))
}

func (m *moduleFactory) OrderModule() {
	fileUsecase := fileUsecases.FileUsecases(m.server.cfg)
	productRepository := productRepositories.ProductRepository(m.server.db, m.server.cfg, fileUsecase)
	repository := orderRepositories.OrderRepository(m.server.db)
	usecase := orderUsecases.OrderUsecases(repository, productRepository)
	handlers := orderHandlers.OrderHandlers(m.server.cfg, usecase)

	router := m.router.Group("/orders")
	_ = router
	_ = handlers
}
