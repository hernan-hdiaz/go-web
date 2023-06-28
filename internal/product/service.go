package product

import (
	"context"
	"fmt"
	"math"
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
	GetTotalPrice(ctx context.Context, productListIds []int) ([]domain.Product, float64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}
func (s *service) GetTotalPrice(ctx context.Context, productListIds []int) ([]domain.Product, float64, error) {
	var productList = []domain.Product{}
	var productQuantity = map[int]int{}
	var totalPrice float64
	for _, id := range productListIds {
		product, err := s.Get(ctx, id)
		if err != nil {
			return []domain.Product{}, 0, err
		}
		if product.IsPublished {
			if productQuantity[product.ID] == 0 && product.Quantity > 0 {
				productQuantity[product.ID] = 1
			} else if productQuantity[product.ID] < product.Quantity {
				productQuantity[product.ID]++
			} else {
				return []domain.Product{}, 0, fmt.Errorf("unavailable quantity for product id: %d", product.ID)
			}
			totalPrice += product.Price
			productList = append(productList, product)
		} else {
			return []domain.Product{}, 0, fmt.Errorf("product not published id: %d", product.ID)
		}
	}
	switch {
	case len(productList) <= 10:
		totalPrice *= 1.21
	case len(productList) > 10 && len(productList) <= 20:
		totalPrice *= 1.17
	default:
		totalPrice *= 1.15
	}

	return productList, roundFloat(totalPrice, 2), nil
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
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
