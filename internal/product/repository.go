package product

import (
	"errors"

	"github.com/hernan-hdiaz/go-web/internal/domain"
	"github.com/hernan-hdiaz/go-web/pkg/store"
)

var (
	ErrNotFound           = errors.New("product not found")
	ErrCreatingProduct    = errors.New("error creating product")
	ErrUpdatingProduct    = errors.New("error updating product")
	ErrAlreadyExists      = errors.New("code_value already exists")
	ErrDateOutOfRange     = errors.New("expiration must be after 01/01/2023")
	ErrPriceOutOfRange    = errors.New("price must be greater than 0")
	ErrQuantityOutOfRange = errors.New("quantity must be greater than 0")
)

type Repository interface {
	GetAll() []domain.Product
	GetByID(id int) (domain.Product, error)
	SearchPriceGt(price float64) []domain.Product
	Create(p domain.Product) (int, error)
	Update(id int, p domain.Product) (domain.Product, error)
	Delete(id int) error
	ValidateCodeValue(codeValue string) bool
}

type repository struct {
	storage store.Store
}

func NewRepository(storage store.Store) Repository {
	return &repository{storage}
}

// retrieves all products
func (r *repository) GetAll() []domain.Product {
	products, err := r.storage.GetAll()
	if err != nil {
		return []domain.Product{}
	}
	return products
}

// search product by ID
func (r *repository) GetByID(id int) (domain.Product, error) {
	product, err := r.storage.GetOne(id)
	if err != nil {
		return domain.Product{}, ErrNotFound
	}
	return product, nil

}

// search for products by price greater than given value
func (r *repository) SearchPriceGt(price float64) []domain.Product {
	var products []domain.Product
	list, err := r.storage.GetAll()
	if err != nil {
		return products
	}
	for _, product := range list {
		if product.Price > price {
			products = append(products, product)
		}
	}
	return products
}

// adds a new product
func (r *repository) Create(p domain.Product) (int, error) {
	if validation := r.ValidateCodeValue(p.CodeValue); !validation {
		return 0, ErrAlreadyExists
	}
	var err error
	p.ID, err = r.storage.AddOne(p)
	if err != nil {
		return 0, ErrCreatingProduct
	}
	return p.ID, nil
}

// validates if the code value already exist on the product list
func (r *repository) ValidateCodeValue(codeValue string) bool {
	list, err := r.storage.GetAll()
	if err != nil {
		return false
	}
	for _, product := range list {
		if product.CodeValue == codeValue {
			return false
		}
	}
	return true
}

// deletes a product
func (r *repository) Delete(id int) error {
	err := r.storage.DeleteOne(id)
	if err != nil {
		return err
	}
	return nil
}

// updates a product
func (r *repository) Update(id int, p domain.Product) (domain.Product, error) {
	err := r.storage.UpdateOne(p)
	if err != nil {
		return domain.Product{}, ErrUpdatingProduct
	}
	return p, nil
}
