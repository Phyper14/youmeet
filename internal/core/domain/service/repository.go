package service

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	CreateService(ctx context.Context, service *Service) error
	GetServiceByID(ctx context.Context, id uuid.UUID) (*Service, error)
	ListServices(ctx context.Context) ([]*Service, error)
}