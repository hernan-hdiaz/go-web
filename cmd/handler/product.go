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

func (p *Product) Update() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Get code_value from path param
		codeValue := c.Param("codeValue")

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

		productUpdated, err := p.productService.Update(c, productRequest, codeValue)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, productUpdated)
	}
}

func (p *Product) Modify() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Get code_value from path param
		codeValue := c.Param("codeValue")

		// //Search warehouse by ID
		// warehouseToUpdate, err := p.productService.Get(c, warehouseToUpdateId)
		// if err != nil {
		// 	web.Error(c, http.StatusNotFound, err.Error())
		// 	return
		// }

		// type productRequest struct {
		// 	ID          int     `json:"id"`
		// 	Name        string  `json:"name"`
		// 	Quantity    int     `json:"quantity"`
		// 	CodeValue   string  `json:"code_value"`
		// 	IsPublished bool    `json:"is_published"`
		// 	Expiration  string  `json:"expiration"`
		// 	Price       float64 `json:"price"`
		// }

		var prodReq domain.ProductRequest
		if err := c.ShouldBindJSON(&prodReq); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			return
		}

		// si no está vacío el Expiration
		// //Parse given date
		// _, err := time.Parse("02/01/2006", prodReq.Expiration)
		// //Check valid format
		// if err != nil {
		// 	c.JSON(http.StatusUnprocessableEntity, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		productModified, err := p.productService.Modify(c, prodReq, codeValue)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, productModified)
	}
}

// func (w *Warehouse) Update() gin.HandlerFunc {
// return func(c *gin.Context) {
//Get ID from path param
// warehouseToUpdateId, err := strconv.Atoi(c.Param("id"))
// if err != nil {
// 	web.Error(c, http.StatusBadRequest, "invalid ID")
// 	return
// }

//Search warehouse by ID
// 	warehouseToUpdate, err := w.warehouseService.Get(c, warehouseToUpdateId)
// 	if err != nil {
// 		web.Error(c, http.StatusNotFound, err.Error())
// 		return
// 	}

// 	//Patch warehouse data
// 	err = c.ShouldBindJSON(&warehouseToUpdate)
// 	if err != nil {
// 		web.Error(c, http.StatusUnprocessableEntity, err.Error())
// 		return
// 	}

// 	//Check id restraint
// 	if warehouseToUpdateId != warehouseToUpdate.ID {
// 		web.Error(c, http.StatusConflict, "cannot modify ID")
// 		return
// 	}

// 	//Check fields restraints
// 	validity := Validation(warehouseToUpdate)
// 	if validity != "" {
// 		web.Error(c, http.StatusUnprocessableEntity, validity)
// 		return
// 	}

// 	//Update warehouse
// 	err = w.warehouseService.Update(c, warehouseToUpdate)
// 	if err != nil {
// 		if errors.Is(err, warehouse.ErrAlreadyExists) || errors.Is(err, warehouse.ErrLocalityID) {
// 			web.Error(c, http.StatusConflict, err.Error())
// 			return
// 		}
// 		web.Error(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	web.Success(c, http.StatusOK, warehouseToUpdate)
// }
// }
