package main

import (
	"net/http"

	"gadget-marketplace/config"
	"gadget-marketplace/handler"
	"gadget-marketplace/pkg/email"
	"gadget-marketplace/repository"
	"gadget-marketplace/router"
	"gadget-marketplace/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.LoadConfig()

	db := config.InitDB(cfg)

	emailSvc := email.NewEmailService(cfg)

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	userService := service.NewUserService(userRepo, emailSvc, cfg)
	productService := service.NewProductService(productRepo, cfg)
	orderService := service.NewOrderService(orderRepo, emailSvc)

	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"app":       "Gadget Marketplace API",
			"framework": "Echo v4",
			"status":    "running",
			"version":   "1.0.0",
		})
	})

	router.SetupRouter(e, cfg, userHandler, productHandler, orderHandler)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
