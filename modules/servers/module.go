package servers

import (
	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/orders/ordersHandlers"
	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/orders/ordersRepositories"
	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/orders/ordersUsecases"

	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/middlewares/middlewaresHandlers"
	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/middlewares/middlewaresRepositories"
	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/middlewares/middlewaresUsecases"

	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/monitor/monitorHandlers"

	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/users/usersHandlers"
	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/users/usersRepositories"

	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/appinfo/appinfoHandlers"
	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/appinfo/appinfoRepositories"
	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/appinfo/appinfoUsecases"

	"github.com/gofiber/fiber/v2"
	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule() IFilesModule
	ProductsModule() IProductsModule
	OrdersModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignOut)

	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecase)

	router := m.r.Group("/appinfo")

	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)

	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)

	router.Delete("/:category_id/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCategory)
}

func (m *moduleFactory) OrdersModule() {
	ordersRepository := ordersRepositories.OrdersRepository(m.s.db)
	ordersUsecase := ordersUsecases.OrdersUsecase(ordersRepository, m.ProductsModule().Repository())
	ordersHandler := ordersHandlers.OrdersHandler(m.s.cfg, ordersUsecase)

	router := m.r.Group("/orders")

	router.Post("/", m.mid.JwtAuth(), ordersHandler.InsertOrder)

	router.Get("/", m.mid.JwtAuth(), m.mid.Authorize(2), ordersHandler.FindOrder)
	router.Get("/:user_id/:order_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), ordersHandler.FindOneOrder)

	router.Patch("/:user_id/:order_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), ordersHandler.UpdateOrder)
}
