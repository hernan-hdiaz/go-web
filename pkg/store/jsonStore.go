package store

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/hernan-hdiaz/go-web/internal/domain"
)

var ErrNotFound = errors.New("product not found")

type Store interface {
	GetAll() ([]domain.Product, error)
	GetOne(id int) (domain.Product, error)
	AddOne(product domain.Product) (int, error)
	UpdateOne(product domain.Product) error
	DeleteOne(id int) error
	saveProducts(products []domain.Product) error
	loadProducts() ([]domain.Product, error)
}

type jsonStore struct {
	pathToFile string
}

// loads products from JSON file
func (s *jsonStore) loadProducts() ([]domain.Product, error) {
	var products []domain.Product
	file, err := os.ReadFile(s.pathToFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(file), &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// saves products to JSON file
func (s *jsonStore) saveProducts(products []domain.Product) error {
	bytes, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return os.WriteFile(s.pathToFile, bytes, 0644)
}

// creates a new product store
func NewStore(path string) Store {
	return &jsonStore{
		pathToFile: path,
	}
}

// retrieves all products
func (s *jsonStore) GetAll() ([]domain.Product, error) {
	products, err := s.loadProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

// search product by id
func (s *jsonStore) GetOne(id int) (domain.Product, error) {
	products, err := s.loadProducts()
	if err != nil {
		return domain.Product{}, err
	}
	for _, product := range products {
		if product.ID == id {
			return product, nil
		}
	}
	return domain.Product{}, ErrNotFound
}

// adds a new product
func (s *jsonStore) AddOne(product domain.Product) (int, error) {
	products, err := s.loadProducts()
	if err != nil {
		return 0, err
	}
	product.ID = len(products) + 1
	products = append(products, product)
	if err = s.saveProducts(products); err != nil {
		return 0, err
	}
	return product.ID, nil
}

// updates a product
func (s *jsonStore) UpdateOne(product domain.Product) error {
	products, err := s.loadProducts()
	if err != nil {
		return err
	}
	for i, p := range products {
		if p.ID == product.ID {
			products[i] = product
			return s.saveProducts(products)
		}
	}
	return ErrNotFound
}

// deletes a product
func (s *jsonStore) DeleteOne(id int) error {
	products, err := s.loadProducts()
	if err != nil {
		return err
	}
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			return s.saveProducts(products)
		}
	}
	return ErrNotFound
}
