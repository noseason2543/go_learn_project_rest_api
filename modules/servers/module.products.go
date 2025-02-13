package servers

import (
	"go_learn_project_rest_api/modules/products/productHandlers"
	"go_learn_project_rest_api/modules/products/productRepositories"
	"go_learn_project_rest_api/modules/products/productUsecases"
)

type IProductsModule interface {
	Init()
	Repository() productRepositories.IProductRepository
	Usecase() productUsecases.IProductUsecases
	Handler() productHandlers.IProductHandler
}

type productsModule struct {
	*moduleFactory
	repository productRepositories.IProductRepository
	usecase    productUsecases.IProductUsecases
	handler    productHandlers.IProductHandler
}

func (m *moduleFactory) ProductModule() IProductsModule {
	repository := productRepositories.ProductRepository(m.server.db, m.server.cfg, m.FilesModule().Usecase())
	usecase := productUsecases.ProductUsecases(repository)
	handlers := productHandlers.ProductHandler(usecase, m.server.cfg, m.FilesModule().Usecase())

	return &productsModule{
		moduleFactory: m,
		repository:    repository,
		usecase:       usecase,
		handler:       handlers,
	}
}

func (p *productsModule) Init() {
	router := p.router.Group("/products")
	router.Post("/addProduct", p.handler.AddProduct, p.mid.JwtAuth(), p.mid.Authorize(2))
	router.Patch("/:product_id", p.handler.UpdateProduct, p.mid.JwtAuth(), p.mid.Authorize(2))
	router.Get("/", p.handler.FindProduct, p.mid.ApiKeyAuth())
	router.Get("/:product_id", p.handler.FindOneProduct, p.mid.ApiKeyAuth())

	router.Delete("/:product_id", p.handler.DeleteProduct, p.mid.JwtAuth(), p.mid.Authorize(2))
}

func (p *productsModule) Repository() productRepositories.IProductRepository { return p.repository }
func (p *productsModule) Usecase() productUsecases.IProductUsecases          { return p.usecase }
func (p *productsModule) Handler() productHandlers.IProductHandler           { return p.handler }
