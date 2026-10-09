package router

import (
	"gadget-marketplace/config"
	"gadget-marketplace/handler"
	"gadget-marketplace/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRouter(
	e *echo.Echo,
	cfg *config.Config,
	userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler,
	orderHandler *handler.OrderHandler,
	aiHandler *handler.AIHandler,
) {
	api := e.Group("/api/v1")

	users := api.Group("/users")
	users.POST("/register", userHandler.Register)
	users.POST("/login", userHandler.Login)

	protectedUsers := users.Group("", middleware.JWTAuthMiddleware(cfg.JWTSecret))
	protectedUsers.GET("/profile", userHandler.GetProfile)
	protectedUsers.POST("/topup", userHandler.TopUp)

	products := api.Group("/products")
	products.POST("/sync", productHandler.SyncProducts)
	products.GET("", productHandler.GetAllProducts)
	products.GET("/:id", productHandler.GetProductByID)

	orders := api.Group("/orders", middleware.JWTAuthMiddleware(cfg.JWTSecret))
	orders.POST("/checkout", orderHandler.Checkout)
	orders.GET("/my-orders", orderHandler.GetUserOrders)

	ai := api.Group("/ai", middleware.JWTAuthMiddleware(cfg.JWTSecret))
	ai.POST("/recommend", aiHandler.Recommend)

	api.POST("/rent-products", orderHandler.Checkout, middleware.JWTAuthMiddleware(cfg.JWTSecret))
	api.GET("/booking-report", orderHandler.GetUserOrders, middleware.JWTAuthMiddleware(cfg.JWTSecret))
}
