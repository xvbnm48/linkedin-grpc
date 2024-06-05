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

func (c Client) GetAllCustomers(ctx context.Context, emailAddress string) ([]models.Customer, error) {
	var customers []models.Customer
	result := c.DB.WithContext(ctx).Where(models.Customer{Email: emailAddress}).Find(&customers)
	return customers, result.Error
}

func (c Client) AddCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	customer.CustomerID = uuid.NewString()
	result := c.DB.WithContext(ctx).Create(&customer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}
	return customer, nil
}

//func (c Client) GetCustomerById(ctx context.Context, ID string) (*models.Customer, error) {
//	customer := &models.Customer{}
//	//result := c.DB.WithContext(ctx).Where("customer_id = ?", ID).First(customer)
//	//result := c.DB.WithContext(ctx).Raw("SELECT * FROM wisdom.customers WHERE customer_id = ?", ID).Scan(customer)
//	result := c.DB.WithContext(ctx).Where(&models.Service{ServiceID: ID}).First(&customer)
//	fmt.Println("result", result)
//	if result.Error != nil {
//		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//			return nil, &dberrors.NotFoundError{Entity: "customer", ID: ID}
//		}
//		return nil, result.Error
//	}
//	//if result.RowsAffected == 0 {
//	//	return nil, &dberrors.NotFoundError{Entity: "customer", ID: ID}
//	//}
//	return customer, nil
//}

func (c Client) GetCustomerById(ctx context.Context, ID string) (*models.Customer, error) {
	customer := &models.Customer{}
	result := c.DB.WithContext(ctx).Where(&models.Customer{CustomerID: ID}).First(&customer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{
				Entity: "customer",
				ID:     ID,
			}
		}
		return nil, result.Error
	}
	return customer, nil
}

func (c Client) UpdateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	var customers []models.Customer
	result := c.DB.WithContext(ctx).
		Model(&customers).
		Clauses(clause.Returning{}).
		Where(&models.Customer{CustomerID: customer.CustomerID}).
		Updates(models.Customer{
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Email:     customer.Email,
			Phone:     customer.Phone,
			Address:   customer.Address,
		})

	fmt.Println("result", result)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{
			Entity: "customer",
			ID:     customer.CustomerID,
		}
	}

	return &customers[0], nil
}
