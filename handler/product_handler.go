package handler

import (
	"net/http"
	"strconv"

	"gadget-marketplace/service"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) SyncProducts(c echo.Context) error {
	products, err := h.productService.SyncFromDummyJSON()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to sync products from 3rd Party API: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Successfully synced product catalog from 3rd Party DummyJSON API",
		"synced_count": len(products),
		"products":     products,
	})
}

func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	category := c.QueryParam("category")
	products, err := h.productService.GetAllProducts(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to retrieve products: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"count":    len(products),
		"products": products,
	})
}

func (h *ProductHandler) GetProductByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid product ID"})
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"product": product,
	})
}
