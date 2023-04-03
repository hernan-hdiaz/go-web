package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var (
	ErrCanNotOpen  = errors.New("can not open file")
	ErrCanNotRead  = errors.New("can not read file")
	ErrCanNotParse = errors.New("can not parse file")
)
var products []Product

func main() {
	err := obtainData()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.GET("/products", getAll)
	router.GET("/products/:id", getById)

	router.Run()
}

func obtainData() error {
	file, err := os.Open("./products.json")
	if err != nil {
		return ErrCanNotOpen
	}
	defer file.Close()

	myDecoder := json.NewDecoder(file)

	if err := myDecoder.Decode(&products); err != nil {
		return ErrCanNotRead
	}
	return nil
}

func getAll(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

func getById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "can not parse")
		return
	}
	for _, product := range products {
		if id == product.ID {
			c.JSON(http.StatusOK, product)
			return
		}
	}
}
