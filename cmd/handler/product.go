package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hernan-hdiaz/go-web/internal/domain"
	"github.com/hernan-hdiaz/go-web/internal/product"
	"github.com/hernan-hdiaz/go-web/pkg/web"
)

var (
	ErrInvalidID    = errors.New("invalid id")
	ErrCanNotParse  = errors.New("can not parse")
	ErrInvalidToken = errors.New("invalid token")
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
			web.Failure(c, http.StatusBadRequest, ErrInvalidID)
			return
		}
		//Search product by ID
		product, err := p.productService.Get(c, id)
		if err != nil {
			web.Failure(c, http.StatusNotFound, err)
			return
		}
		//Return found product
		web.Success(c, http.StatusOK, product)
	}
}

func (p *Product) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products := p.productService.GetAll(c)
		//Return products
		web.Success(c, http.StatusOK, products)
	}
}

func (p *Product) GetTotalPrice() gin.HandlerFunc {
	return func(c *gin.Context) {
		productList := c.Query("list")
		productList, _ = strings.CutPrefix(productList, "[")
		productList, _ = strings.CutSuffix(productList, "]")
		productListIds := strings.Split(productList, ",")

		var convertedProductListIds = []int{}
		for _, p := range productListIds {
			productId, err := strconv.Atoi(p)
			if err != nil {
				web.Failure(c, http.StatusBadRequest, ErrInvalidID)
				return
			}
			convertedProductListIds = append(convertedProductListIds, productId)
		}
		completeProductList, totalPrice, err := p.productService.GetTotalPrice(c, convertedProductListIds)
		if err != nil {
			web.Failure(c, http.StatusBadRequest, err)
			return
		}

		type consumerPrice struct {
			Products   []domain.Product `json:"products"`
			TotalPrice float64          `json:"total_price"`
		}

		var response = consumerPrice{Products: completeProductList, TotalPrice: totalPrice}

		//Return products
		web.Success(c, http.StatusOK, response)
	}
}

func (p *Product) SearchByPriceGt() gin.HandlerFunc {
	return func(c *gin.Context) {
		priceGt, err := strconv.ParseFloat(c.Query("priceGt"), 64)
		if err != nil {
			web.Failure(c, http.StatusBadRequest, ErrCanNotParse)
			return
		}

		productsByPriceGt, err := p.productService.SearchByPriceGt(c, priceGt)
		if err != nil {
			web.Failure(c, http.StatusNotFound, err)
			return
		}
		web.Success(c, http.StatusOK, productsByPriceGt)
	}
}

func (p *Product) Save() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productRequest domain.Product
		if err := c.ShouldBindJSON(&productRequest); err != nil {
			web.Failure(c, http.StatusUnprocessableEntity, err)
			return
		}

		//Parse given date
		_, err := time.Parse("02/01/2006", productRequest.Expiration)
		//Check valid format
		if err != nil {
			web.Failure(c, http.StatusUnprocessableEntity, err)
			return
		}

		productRequest.ID, err = p.productService.Save(c, productRequest)
		if err != nil {
			web.Failure(c, http.StatusConflict, err)
			return
		}
		web.Success(c, http.StatusCreated, productRequest)
	}
}

func (p *Product) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get ID from path param
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Failure(c, http.StatusBadRequest, ErrInvalidID)
			return
		}

		var productRequest domain.ProductRequest
		if err := c.ShouldBindJSON(&productRequest); err != nil {
			web.Failure(c, http.StatusUnprocessableEntity, err)
			return
		}
		if productRequest.Expiration != "" {
			//Parse given date
			_, err = time.Parse("02/01/2006", productRequest.Expiration)
			//Check valid format
			if err != nil {
				web.Failure(c, http.StatusUnprocessableEntity, err)
				return
			}
		}
		productUpdated, err := p.productService.Update(c, productRequest, id)
		if err != nil {
			web.Failure(c, http.StatusNotFound, err)
			return
		}
		web.Success(c, http.StatusCreated, productUpdated)
	}
}

func (p *Product) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get ID from path param
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Failure(c, http.StatusBadRequest, ErrInvalidID)
			return
		}
		//Search product by codeValue
		err = p.productService.Delete(c, id)
		if err != nil {
			web.Failure(c, http.StatusNotFound, err)
			return
		}
		web.Success(c, http.StatusNoContent, nil)
	}
}
