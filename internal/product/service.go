package product

import (
	"context"
	"time"

	"github.com/hernan-hdiaz/go-web/internal/domain"
)

type Service interface {
	Get(ctx context.Context, id int) (domain.Product, error)
	GetAll(ctx context.Context) []domain.Product
	SearchByPriceGt(ctx context.Context, priceGt float64) ([]domain.Product, error)
	Save(ctx context.Context, productRequest domain.Product) (int, error)
	Update(ctx context.Context, productRequest domain.ProductRequest, id int) (domain.Product, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Get(ctx context.Context, id int) (domain.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (s *service) GetAll(ctx context.Context) []domain.Product {
	products := s.repo.GetAll()
	return products
}

func (s *service) SearchByPriceGt(ctx context.Context, priceGt float64) ([]domain.Product, error) {
	products := s.repo.SearchPriceGt(priceGt)
	return products, nil
}

func (s *service) Save(ctx context.Context, productRequest domain.Product) (int, error) {
	date, _ := time.Parse("02/01/2006", productRequest.Expiration)
	//Set minimum date
	minimum_date, _ := time.Parse("02/01/2006", "01/01/2023")
	//Check date restraints
	if date.Before(minimum_date) {
		return 0, ErrDateOutOfRange
	}
	if productRequest.Price <= 0 {
		return 0, ErrPriceOutOfRange
	}
	if productRequest.Quantity <= 0 {
		return 0, ErrQuantityOutOfRange
	}

	productID, err := s.repo.Create(productRequest)
	if err != nil {
		return 0, err
	}
	return productID, nil
}

func (s *service) Update(ctx context.Context, productRequest domain.ProductRequest, id int) (domain.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return domain.Product{}, err
	}
	if productRequest.Name != "" {
		product.Name = productRequest.Name
	}
	if productRequest.CodeValue != "" && productRequest.CodeValue != product.CodeValue {
		validation := s.repo.ValidateCodeValue(productRequest.CodeValue)
		if !validation {
			return domain.Product{}, ErrAlreadyExists
		}

		product.CodeValue = productRequest.CodeValue
	}
	if productRequest.Expiration != "" {
		date, _ := time.Parse("02/01/2006", productRequest.Expiration)
		//Set minimum date
		minimum_date, _ := time.Parse("02/01/2006", "01/01/2023")
		//Check date restraints
		if date.Before(minimum_date) {
			return domain.Product{}, ErrDateOutOfRange
		}
		product.Expiration = productRequest.Expiration
	}
	if productRequest.Quantity > 0 {
		product.Quantity = productRequest.Quantity
	}
	if productRequest.Price > 0 {
		product.Price = productRequest.Price
	}
	if productRequest.IsPublished != nil {
		product.IsPublished = *productRequest.IsPublished
	}
	product, err = s.repo.Update(id, product)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(id)
	return err
}
