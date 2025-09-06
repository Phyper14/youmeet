package repositories

import (
	"context"
	"github.com/google/uuid"
	"youmeet/internal/core/domain/service"
)

type ServiceRepository struct {
	db DBClient
}

func NewServiceRepository(db DBClient) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) CreateService(ctx context.Context, svc *service.Service) error {
	return r.db.Create(svc)
}

func (r *ServiceRepository) GetServiceByID(ctx context.Context, id uuid.UUID) (*service.Service, error) {
	var svc service.Service
	err := r.db.First(&svc, "id = ?", id)
	return &svc, err
}

func (r *ServiceRepository) ListServices(ctx context.Context) ([]*service.Service, error) {
	var services []*service.Service
	err := r.db.Find(&services)
	return services, err
}