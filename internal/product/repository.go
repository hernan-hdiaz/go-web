package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hernan-hdiaz/go-web/internal/domain"
)

var (
	ErrCanNotOpen         = errors.New("can not open file")
	ErrCanNotRead         = errors.New("can not read file")
	ErrNotFound           = errors.New("product not found")
	ErrAlreadyExists      = errors.New("code_value already exists")
	ErrDateOutOfRange     = errors.New("expiration must be after 01/01/2023")
	ErrCodeValueMissmatch = errors.New("codeValues missmatch")
	products              = []domain.Product{}
)

// Repository encapsulates the storage of a product.
type Repository interface {
	Get(ctx context.Context, id int) (domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	SearchByCodeValue(ctx context.Context, codeValue string) (domain.Product, int)
	SearchByPriceGt(ctx context.Context, priceGt float64) ([]domain.Product, error)
	Save(ctx context.Context, productRequest domain.Product) int
	Exists(ctx context.Context, codeValue string) bool
	Update(ctx context.Context, product domain.Product, index int)
}

type repository struct {
	db *os.File
}

func NewRepository() (Repository, error) {
	file, err := os.Open("./products.json")
	if err != nil {
		return &repository{}, ErrCanNotOpen
	}
	defer file.Close()

	myDecoder := json.NewDecoder(file)

	if err := myDecoder.Decode(&products); err != nil {
		return &repository{}, ErrCanNotRead
	}
	return &repository{
		db: file,
	}, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Product, error) {
	for _, product := range products {
		if id == product.ID {
			return product, nil
		}
	}
	return domain.Product{}, ErrNotFound
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Product, error) {
	return products, nil
}

func (r *repository) SearchByPriceGt(ctx context.Context, priceGt float64) ([]domain.Product, error) {
	var productsByPriceGt []domain.Product

	for _, product := range products {
		if priceGt < product.Price {
			productsByPriceGt = append(productsByPriceGt, product)
		}
	}
	if len(productsByPriceGt) == 0 {
		return []domain.Product{}, fmt.Errorf("not found products with price greater than %.2f", priceGt)
	}
	return productsByPriceGt, nil
}

func (r *repository) Save(ctx context.Context, productRequest domain.Product) int {
	lastID := len(products) + 1
	productRequest.ID = lastID
	products = append(products, productRequest)
	return productRequest.ID
}

func (r *repository) Exists(ctx context.Context, codeValue string) bool {
	for _, p := range products {
		if p.CodeValue == codeValue {
			return true
		}
	}

	return false
}

func (r *repository) SearchByCodeValue(ctx context.Context, codeValue string) (domain.Product, int) {
	for i, p := range products {
		if p.CodeValue == codeValue {
			return p, i
		}
	}

	return domain.Product{}, 0
}

func (r *repository) Update(ctx context.Context, productRequest domain.Product, index int) {
	products = append(products[:index], products[index+1:]...)
	products = append(products, productRequest)
}

// Update
// func main() {
// 	var strSlice = []string{"India", "Canada", "Japan", "Germany", "Italy"}
// 	fmt.Println(strSlice)

// 	strSlice = RemoveIndex(strSlice, 3)
// 	fmt.Println(strSlice)
// }

// func RemoveIndex(s []string, index int) []string {
// 	return append(s[:index], s[index+1:]...)
// }