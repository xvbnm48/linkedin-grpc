package database

import (
	"context"
	"fmt"
	"github.com/xvbnm48/linkedin-grpc/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DatabaseClient interface {
	Ready() bool

	GetAllCustomers(ctx context.Context, emailAddress string) ([]models.Customer, error)
	AddCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error)
	GetCustomerById(ctx context.Context, ID string) (*models.Customer, error)
	UpdateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error)

	GetAllProducts(ctx context.Context, vendorId string) ([]models.Product, error)
	AddProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	GetProductById(ctx context.Context, ID string) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)

	GetAllService(ctx context.Context) ([]models.Service, error)
	AddService(ctx context.Context, service *models.Service) (*models.Service, error)
	GetServiceByID(ctx context.Context, ID string) (*models.Service, error)
	UpdateService(ctx context.Context, service *models.Service) (*models.Service, error)

	GetAllVendors(ctx context.Context) ([]models.Vendor, error)
	AddVendor(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error)
	GetVendorByID(ctx context.Context, id string) (*models.Vendor, error)
}

type Client struct {
	DB *gorm.DB
}

func NewDatabaseClient() (DatabaseClient, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		"localhost",
		"postgres",
		"fariz",
		"linkedid-grpc",
		5432,
		"disable",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "wisdom.", // table name prefix
		},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		QueryFields: true,
	})
	if err != nil {
		return nil, err
	}
	client := Client{DB: db}

	return client, nil
}

func (c Client) Ready() bool {
	var ready string
	tx := c.DB.Raw("SELECT 1 as ready").Scan(&ready)
	if tx.Error != nil {
		return false
	}
	if ready == "1" {
		return true
	}
	return false
}
