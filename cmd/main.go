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

// Ejercicio 1: Dominios
// Es momento de organizar nuestra API, seguiremos el Diseño orientado a paquetes –
// “Package Oriented Design”–. Para empezar vamos a dividir nuestro proyecto en dominios,
// por ahora, solo tenemos uno: product.
// Luego, tendremos que refactorizar nuestro código a la estructura por dominios, repository,
// service y controller (o también handler). Recordemos lo visto en clase:
// Repository: abstrae el acceso a los datos.
// Service: contiene la lógica de negocio de la API, maneja conexiones externas.
// Controller: toma las peticiones del cliente, valida las entradas y retorna las respuestas.

// Ejercicio 1: Método PUT
// Añadir el método PUT a nuestra API, recordemos que crea o reemplaza un recurso en su totalidad
// con el contenido en la request. Tené en cuenta validar los campos que se envían, como hiciste con
// el método POST. Seguimos aplicando los cambios sobre la lista cargada en memoria.

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/hernan-hdiaz/go-web/cmd/handler"
	"github.com/hernan-hdiaz/go-web/internal/domain"
	"github.com/hernan-hdiaz/go-web/internal/product"
)

var (
	ErrCanNotOpen     = errors.New("can not open file")
	ErrCanNotRead     = errors.New("can not read file")
	ErrCanNotParse    = errors.New("can not parse file")
	ErrAlreadyExists  = errors.New("code_value already exists")
	ErrDateOutOfRange = errors.New("expiration must be after 01/01/2023")
	Products          = []domain.Product{}
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	repo, err := product.NewRepository()
	if err != nil {
		panic(err)
	}
	service := product.NewService(repo)
	handler := handler.NewProductHandler(service)

	router.GET("/products", handler.GetAll())
	router.GET("/products/:id", handler.Get())
	router.GET("/products/search", handler.SearchByPriceGt())
	router.POST("/products", handler.Save())
	router.PUT("/products/:codeValue", handler.Update())
	router.PATCH("/products/:codeValue", handler.Modify())
	router.DELETE("/products/:codeValue", handler.Delete())

	router.Run()
}
