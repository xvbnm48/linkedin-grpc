package database

import (
	"context"
	"github.com/xvbnm48/linkedin-grpc/internal/models"
)

func (c Client) GetAllService(ctx context.Context) ([]models.Service, error) {
	var service []models.Service
	result := c.DB.WithContext(ctx).Find(&service)
	return service, result.Error
}
