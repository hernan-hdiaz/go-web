package product

import (
	"context"
	"time"

	"github.com/hernan-hdiaz/go-web/internal/domain"
)

type Service interface {
	Get(ctx context.Context, id int) (domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	SearchByPriceGt(ctx context.Context, priceGt float64) ([]domain.Product, error)
	Save(ctx context.Context, productRequest domain.Product) (int, error)
	Update(ctx context.Context, productRequest domain.Product, codeValue string) (domain.Product, error)
	Modify(ctx context.Context, productRequest domain.ProductRequest, codeValue string) (domain.Product, error)
	Delete(ctx context.Context, codeValue string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Get(ctx context.Context, id int) (domain.Product, error) {
	product, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (s *service) GetAll(ctx context.Context) ([]domain.Product, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return []domain.Product{}, err
	}
	return products, nil
}

func (s *service) SearchByPriceGt(ctx context.Context, priceGt float64) ([]domain.Product, error) {
	products, err := s.repo.SearchByPriceGt(ctx, priceGt)
	if err != nil {
		return []domain.Product{}, err
	}
	return products, nil
}

func (s *service) Save(ctx context.Context, productRequest domain.Product) (int, error) {
	//Check code_value
	if s.repo.Exists(ctx, productRequest.CodeValue) {
		return 0, ErrAlreadyExists
	}

	date, _ := time.Parse("02/01/2006", productRequest.Expiration)
	//Set minimum date
	minimum_date, _ := time.Parse("02/01/2006", "01/01/2023")
	//Check date restraints
	if date.Before(minimum_date) {
		return 0, ErrDateOutOfRange
	}

	productID := s.repo.Save(ctx, productRequest)
	return productID, nil
}

func (s *service) Update(ctx context.Context, productRequest domain.Product, codeValue string) (domain.Product, error) {

	date, _ := time.Parse("02/01/2006", productRequest.Expiration)
	//Set minimum date
	minimum_date, _ := time.Parse("02/01/2006", "01/01/2023")
	//Check date restraints
	if date.Before(minimum_date) {
		return domain.Product{}, ErrDateOutOfRange
	}

	//Check code_values
	if codeValue != productRequest.CodeValue {
		if s.repo.Exists(ctx, productRequest.CodeValue) {
			return domain.Product{}, ErrAlreadyExists
		} else if !s.repo.Exists(ctx, codeValue) {
			return domain.Product{}, ErrCodeValueMissmatch
		}
	}

	if s.repo.Exists(ctx, codeValue) {
		product, index, err := s.repo.SearchByCodeValue(ctx, codeValue)
		if err != nil {
			return domain.Product{}, err
		}
		productRequest.ID = product.ID
		s.repo.Update(ctx, productRequest, index)
	} else {
		productRequest.ID = s.repo.Save(ctx, productRequest)
	}

	return productRequest, nil
}

func (s *service) Modify(ctx context.Context, productRequest domain.ProductRequest, codeValue string) (domain.Product, error) {

	if productToModify, index, err := s.repo.SearchByCodeValue(ctx, codeValue); err == nil {
		//Check expiration restraints
		if productRequest.Expiration != "" {
			date, err := time.Parse("02/01/2006", productRequest.Expiration)
			if err != nil {
				return domain.Product{}, err
			}
			//Set minimum date
			minimum_date, _ := time.Parse("02/01/2006", "01/01/2023")
			//Check date restraints
			if date.Before(minimum_date) {
				return domain.Product{}, ErrDateOutOfRange
			}
		}
		//Check code_values restraints
		if productRequest.CodeValue != "" {
			if codeValue != productRequest.CodeValue {
				if s.repo.Exists(ctx, productRequest.CodeValue) {
					return domain.Product{}, ErrAlreadyExists
				} else if !s.repo.Exists(ctx, codeValue) {
					return domain.Product{}, ErrCodeValueMissmatch
				}
			}
		}
		//Check quantity restraints
		if productRequest.Quantity < 0 {
			return domain.Product{}, ErrQuantityOutOfRange
		}
		//Check price restraints
		if productRequest.Price < 0 {
			return domain.Product{}, ErrPriceOutOfRange
		}

		// Patch product data
		if productRequest.Name != "" {
			productToModify.Name = productRequest.Name
		}
		if productRequest.Quantity != 0 {
			productToModify.Quantity = productRequest.Quantity
		}
		if productRequest.CodeValue != "" {
			productToModify.CodeValue = productRequest.CodeValue
		}
		if productRequest.IsPublished != nil {
			productToModify.IsPublished = *productRequest.IsPublished
		}
		if productRequest.Expiration != "" {
			productToModify.Expiration = productRequest.Expiration
		}
		if productRequest.Price != 0 {
			productToModify.Price = productRequest.Price
		}

		s.repo.Update(ctx, productToModify, index)
		return productToModify, nil
	} else {
		return domain.Product{}, ErrNotFound
	}
}

func (s *service) Delete(ctx context.Context, codeValue string) error {
	err := s.repo.Delete(ctx, codeValue)
	return err
}
