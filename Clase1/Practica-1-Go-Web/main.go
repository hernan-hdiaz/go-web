package main

// Ejercicio 1 : Iniciando el proyecto
// Debemos crear un repositorio en github.com para poder subir nuestros avances.
// Este repositorio es el que vamos a utilizar para llevar lo que realicemos durante las
// distintas prácticas de Go Web.
// Primero debemos clonar el repositorio creado, luego iniciar nuestro proyecto de go con con el
// comando go mod init.
// El siguiente paso será crear un archivo main.go donde deberán cargar en una slice, desde un
// archivo JSON, los datos de productos. Esta slice se debe cargar cada vez que se inicie la API
// para realizar las distintas consultas.

// Ejercicio 2 : Creando un servidor
// Vamos a levantar un servidor utilizando el paquete gin en el puerto 8080.
// Para probar nuestros endpoints haremos uso de postman.
// Crear una ruta /ping que debe respondernos con un string que contenga pong con el status 200 OK.
// Crear una ruta /products que nos devuelva la lista de todos los productos en la slice.
// Crear una ruta /products/:id que nos devuelva un producto por su id.
// Crear una ruta /products/search que nos permita buscar por parámetro los productos cuyo
// precio sean mayor a un valor priceGt.

// Ejercicio 1: Añadir un producto
// En esta ocasión vamos a añadir un producto al slice cargado en memoria.
// Dentro de la ruta /products añadimos el método POST, al cual vamos a enviar en el cuerpo
// de la request el nuevo producto. El mismo tiene ciertas restricciones, conozcámoslas:
// No es necesario pasar el Id, al momento de añadirlo se debe inferir del estado de la lista
// de productos, verificando que no se repitan ya que debe ser un campo único.
// Ningún dato puede estar vacío, exceptuando is_published (vacío indica un valor false).
// El campo code_value debe ser único para cada producto.
// Los tipos de datos deben coincidir con los definidos en el planteo del problema.
// La fecha de vencimiento debe tener el formato: XX/XX/XXXX, además debemos verificar que día,
// mes y año sean valores válidos.
// Recordá: si una consulta está mal formulada por parte del cliente, el status code
// cae en los 4XX.

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	CodeValue   string  `json:"code_value" binding:"required"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

var (
	ErrCanNotOpen     = errors.New("can not open file")
	ErrCanNotRead     = errors.New("can not read file")
	ErrCanNotParse    = errors.New("can not parse file")
	ErrAlreadyExists  = errors.New("code_value already exists")
	ErrDateOutOfRange = errors.New("expiration must be after 01/01/2023")
	products          = []Product{}
)

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
	router.GET("/products/search", searchByPriceGt)
	router.POST("/products", save)

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

func searchByPriceGt(c *gin.Context) {
	price, err := strconv.ParseFloat(c.Query("priceGt"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "can not parse")
		return
	}

	var productsByPriceGt []Product

	for _, product := range products {
		if price < product.Price {
			productsByPriceGt = append(productsByPriceGt, product)
		}
	}

	c.JSON(http.StatusOK, productsByPriceGt)
}

func save(c *gin.Context) {
	var productRequest Product
	if err := c.ShouldBindJSON(&productRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Check code_value
	for _, p := range products {
		if p.CodeValue == productRequest.CodeValue {
			c.JSON(http.StatusConflict, gin.H{
				"error": ErrAlreadyExists.Error(),
			})
			return
		}
	}

	//Parse given date
	date, err := time.Parse("02/01/2006", productRequest.Expiration)
	//Check valid format
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	//Set minimum date
	minimum_date, _ := time.Parse("02/01/2006", "01/01/2023")
	//Check date restraints
	if date.Before(minimum_date) {
		c.JSON(http.StatusConflict, gin.H{
			"error": ErrDateOutOfRange.Error(),
		})
		return
	}

	//Save product
	lastID := len(products) + 1
	productRequest.ID = lastID
	products = append(products, productRequest)
	c.JSON(http.StatusCreated, productRequest)
}
