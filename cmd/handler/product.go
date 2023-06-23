package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hernan-hdiaz/go-web/internal/domain"
	"github.com/hernan-hdiaz/go-web/internal/product"
)

type Product struct {
	productService product.Service
}

func NewProductHandler(p product.Service) *Product {
	return &Product{
		productService: p,
	}
}

func (p *Product) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get ID from path param
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id must be int",
			})
			return
		}
		//Search product by ID
		product, err := p.productService.Get(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		//Return found product
		c.JSON(http.StatusOK, product)
	}
}

func (p *Product) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := p.productService.GetAll(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		//Return products
		c.JSON(http.StatusOK, products)
	}
}

func (p *Product) SearchByPriceGt() gin.HandlerFunc {
	return func(c *gin.Context) {
		priceGt, err := strconv.ParseFloat(c.Query("priceGt"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "can not parse",
			})
			return
		}

		productsByPriceGt, err := p.productService.SearchByPriceGt(c, priceGt)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, productsByPriceGt)
	}
}

func (p *Product) Save() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productRequest domain.Product
		if err := c.ShouldBindJSON(&productRequest); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			return
		}

		//Parse given date
		_, err := time.Parse("02/01/2006", productRequest.Expiration)
		//Check valid format
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			return
		}

		productRequest.ID, err = p.productService.Save(c, productRequest)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, productRequest)
	}
}
