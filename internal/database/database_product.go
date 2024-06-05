package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/xvbnm48/linkedin-grpc/internal/dberrors"
	"github.com/xvbnm48/linkedin-grpc/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c Client) GetAllProducts(ctx context.Context, vendorId string) ([]models.Product, error) {
	var products []models.Product
	result := c.DB.WithContext(ctx).Where(models.Product{VendorID: vendorId}).Find(&products)
	return products, result.Error
}

func (c Client) AddProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	product.ProductID = uuid.NewString()
	//product.VendorID = uuid.NewString()
	fmt.Println(product)
	result := c.DB.WithContext(ctx).Create(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}

	return product, nil
}

func (c Client) GetProductById(ctx context.Context, ID string) (*models.Product, error) {
	product := &models.Product{}
	result := c.DB.WithContext(ctx).Where(&models.Product{ProductID: ID}).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "product", ID: ID}
		}
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{Entity: "product", ID: ID}
	}

	return product, nil
}

func (c Client) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	var products []models.Product
	result := c.DB.WithContext(ctx).
		Model(&products).
		Clauses(clause.Returning{}).
		Where(&models.Product{ProductID: product.ProductID}).
		Updates(models.Product{
			ProductID: product.ProductID,
			Name:      product.Name,
			Price:     product.Price,
			VendorID:  product.VendorID,
		})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{Entity: "product", ID: product.ProductID}
	}

	return &products[0], nil
}
