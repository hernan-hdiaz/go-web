package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hernan-hdiaz/go-web/cmd/handler"
	"github.com/hernan-hdiaz/go-web/internal/domain"
	"github.com/hernan-hdiaz/go-web/internal/product"
	"github.com/hernan-hdiaz/go-web/pkg/store"
	"github.com/joho/godotenv"
)

var (
	Products = []domain.Product{}
)

func main() {
	if err := godotenv.Load("./cmd/server/.env"); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	storage := store.NewStore("./products.json")
	repo := product.NewRepository(storage)
	service := product.NewService(repo)
	handler := handler.NewProductHandler(service)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.GET("/products", handler.GetAll())
	router.GET("/products/:id", handler.Get())
	router.GET("/products/consumer_price", handler.GetTotalPrice())
	router.GET("/products/search", handler.SearchByPriceGt())
	router.Use(TokenAuthMiddleware())
	router.POST("/products", handler.Save())
	router.PUT("/products/:id", handler.Update())
	router.DELETE("/products/:id", handler.Delete())

	router.Run()
}

func TokenAuthMiddleware() gin.HandlerFunc {
	requiredToken := os.Getenv("TOKEN")

	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != requiredToken {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API token"})
			return
		}
		c.Next()
	}
}
