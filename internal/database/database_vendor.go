package database

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/xvbnm48/linkedin-grpc/internal/dberrors"
	"github.com/xvbnm48/linkedin-grpc/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c Client) GetAllVendors(ctx context.Context) ([]models.Vendor, error) {
	var vendor []models.Vendor
	result := c.DB.WithContext(ctx).Find(&vendor)
	return vendor, result.Error

}

func (c Client) AddVendor(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error) {
	vendor.VendorID = uuid.NewString()
	result := c.DB.WithContext(ctx).Create(&vendor)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}
	return vendor, nil
}

func (c Client) GetVendorByID(ctx context.Context, id string) (*models.Vendor, error) {
	vendor := &models.Vendor{}
	// this is a get by id vendor raw query
	result := c.DB.WithContext(ctx).Where(&models.Vendor{VendorID: id}).First(&vendor)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "vendor", ID: id}
		}
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{Entity: "vendor", ID: id}
	}

	return vendor, nil
}

func (c Client) UpdateVendor(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error) {
	var vendors []models.Vendor
	result := c.DB.WithContext(ctx).
		Model(&vendors).
		Clauses(clause.Returning{}).
		Where(&models.Vendor{VendorID: vendor.VendorID}).
		Updates(models.Vendor{
			VendorID: vendor.VendorID,
			Name:     vendor.Name,
			Contact:  vendor.Contact,
			Phone:    vendor.Phone,
			Email:    vendor.Email,
			Address:  vendor.Address,
		})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{
			Entity: "Vendor",
			ID:     vendor.VendorID,
		}
	}
	return &vendors[0], nil
}
