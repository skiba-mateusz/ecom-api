package service

import (
	"context"
	"github.com/skiba-mateusz/ecom-api/internal/app/domain"
	"github.com/skiba-mateusz/ecom-api/internal/infra/persistence/postgres/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetProduct(t *testing.T) {
	t.Run("should_return_product_with_category", func(t *testing.T) {
		mockProductRepo := new(repository.MockProductRepository)
		mockCategoryRepo := new(repository.MockCategoryRepository)
		productServ := NewProductService(mockProductRepo, mockCategoryRepo)

		productId := int64(123)
		categoryId := int64(456)

		mockProduct := &domain.Product{
			BaseProduct: domain.BaseProduct{
				Id:         productId,
				Name:       "Mock Product",
				CategoryId: categoryId,
			},
		}

		mockCategory := &domain.Category{
			Id:   categoryId,
			Name: "Mock Category",
		}

		mockProductRepo.On("GetById", mock.Anything, productId).Return(mockProduct, nil)
		mockCategoryRepo.On("GetById", mock.Anything, categoryId).Return(mockCategory, nil)

		ctx := context.Background()
		result, err := productServ.GetById(ctx, productId)

		assert.NoError(t, err)

		assert.NotNil(t, result)
		assert.Equal(t, mockProduct.Id, result.Id)
		assert.Equal(t, mockProduct.Name, result.Name)
		assert.Equal(t, mockProduct.CategoryId, result.CategoryId)
		assert.NotNil(t, result.Category)
		assert.Equal(t, mockCategory.Id, result.Category.Id)
		assert.Equal(t, mockCategory.Name, result.Category.Name)

		mockProductRepo.AssertExpectations(t)
		mockCategoryRepo.AssertExpectations(t)
	})

	t.Run("should_return_product_with_optional_fields", func(t *testing.T) {
		mockProductRepo := new(repository.MockProductRepository)
		mockCategoryRepo := new(repository.MockCategoryRepository)
		productServ := NewProductService(mockProductRepo, mockCategoryRepo)

		productId := int64(123)
		categoryId := int64(456)

		description := "Mock description"
		salePrice := 99.99

		mockProduct := &domain.Product{
			BaseProduct: domain.BaseProduct{
				Id:         productId,
				Name:       "Mock Product",
				CategoryId: categoryId,
				SalePrice:  &salePrice,
			},
			Description: &description,
		}

		mockCategory := &domain.Category{
			Id:   categoryId,
			Name: "Mock Category",
		}

		mockProductRepo.On("GetById", mock.Anything, productId).Return(mockProduct, nil)
		mockCategoryRepo.On("GetById", mock.Anything, categoryId).Return(mockCategory, nil)

		ctx := context.Background()
		result, err := productServ.GetById(ctx, productId)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockProduct.Id, result.Id)
		assert.Equal(t, mockProduct.Name, result.Name)
		assert.Equal(t, mockProduct.SalePrice, result.SalePrice)
		assert.Equal(t, mockProduct.Description, result.Description)
		assert.Equal(t, mockProduct.CategoryId, result.Category.Id)
		assert.NotNil(t, result.Category)
		assert.Equal(t, mockProduct.Category.Id, result.Category.Id)
		assert.Equal(t, mockProduct.Category.Name, result.Category.Name)

		mockProductRepo.AssertExpectations(t)
		mockCategoryRepo.AssertExpectations(t)
	})

	t.Run("should_return_error_when_product_not_found", func(t *testing.T) {
		mockProductRepo := new(repository.MockProductRepository)
		mockCategoryRepo := new(repository.MockCategoryRepository)
		productServ := NewProductService(mockProductRepo, mockCategoryRepo)

		productId := int64(123)

		mockProductRepo.On("GetById", mock.Anything, productId).Return(nil, domain.ErrNotFound)

		ctx := context.Background()
		result, err := productServ.GetById(ctx, productId)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrNotFound, err)

		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should_return_error_when_category_not_found", func(t *testing.T) {
		mockProductRepo := new(repository.MockProductRepository)
		mockCategoryRepo := new(repository.MockCategoryRepository)
		productServ := NewProductService(mockProductRepo, mockCategoryRepo)

		productId := int64(123)
		categoryId := int64(999)

		mockProduct := &domain.Product{
			BaseProduct: domain.BaseProduct{
				Id:         productId,
				Name:       "Mock Product",
				CategoryId: categoryId,
			},
		}

		mockProductRepo.On("GetById", mock.Anything, productId).Return(mockProduct, nil)
		mockCategoryRepo.On("GetById", mock.Anything, categoryId).Return(nil, domain.ErrNotFound)

		ctx := context.Background()
		result, err := productServ.GetById(ctx, productId)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrNotFound, err)

		mockProductRepo.AssertExpectations(t)
		mockCategoryRepo.AssertExpectations(t)
	})
}

func TestCreateProduct(t *testing.T) {
	t.Run("should_create_product", func(t *testing.T) {
		mockProductRepo := new(repository.MockProductRepository)
		productServ := NewProductService(mockProductRepo, nil)

		productId := int64(123)
		categoryId := int64(456)
		branId := int64(789)

		slug := "mock-product"
		description := "Mock description"

		mockProduct := &domain.Product{
			BaseProduct: domain.BaseProduct{
				Name:       "Mock Product",
				Price:      120,
				SalePrice:  nil,
				CategoryId: categoryId,
				Stock:      10,
				BrandId:    branId,
			},
			Description: &description,
		}

		mockProductRepo.On("SlugExists", mock.Anything, slug).Return(false, nil)
		mockProductRepo.On("Create", mock.Anything, mockProduct).Run(func(args mock.Arguments) {
			product := args.Get(1).(*domain.Product)
			product.Id = productId
		}).Return(nil)

		ctx := context.Background()
		err := productServ.Create(ctx, mockProduct)

		assert.NoError(t, err)

		assert.Equal(t, mockProduct.Id, productId)
		assert.Equal(t, mockProduct.Name, mockProduct.Name)
		assert.Equal(t, mockProduct.Slug, slug)
		assert.Nil(t, mockProduct.SalePrice)
		assert.Equal(t, mockProduct.CategoryId, mockProduct.CategoryId)
		assert.Equal(t, mockProduct.Stock, mockProduct.Stock)
		assert.Equal(t, mockProduct.BrandId, mockProduct.BrandId)
		assert.Equal(t, mockProduct.Description, mockProduct.Description)

		mockProductRepo.AssertExpectations(t)
	})
}
