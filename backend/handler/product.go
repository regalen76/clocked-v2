package handler

import (
	"reonify/clocked/database"
	"reonify/clocked/model"

	"github.com/gofiber/fiber/v2"
)

// GetAllProducts query all products
func GetAllProducts(c *fiber.Ctx) error {
	db := database.DB
	var products []model.Product
	db.Find(&products)
	return c.JSON(fiber.Map{"status": "success", "message": "All products", "data": products})
}

// GetProduct query product
func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var product model.Product
    db.Find(&product, id)
    if product.Title == "" {
        return fiber.NewError(fiber.StatusNotFound, "No product found with ID")
    }
    return c.JSON(fiber.Map{"status": "success", "message": "Product found", "data": product})
}

// CreateProduct new product
func CreateProduct(c *fiber.Ctx) error {
	db := database.DB
    product := new(model.Product)
    if err := c.BodyParser(product); err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Couldn't create product")
    }
    db.Create(&product)
    return c.JSON(fiber.Map{"status": "success", "message": "Created product", "data": product})
}

// DeleteProduct delete product
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var product model.Product
    db.First(&product, id)
    if product.Title == "" {
        return fiber.NewError(fiber.StatusNotFound, "No product found with ID")
    }
    db.Delete(&product)
    return c.JSON(fiber.Map{"status": "success", "message": "Product successfully deleted", "data": nil})
}
