package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type CompanyRepository interface {
	CreateCompany(ctx context.Context, company *Company) error
	GetCompanyByID(ctx context.Context, id uuid.UUID) (*Company, error)
	GetCompanyByUserID(ctx context.Context, userID uuid.UUID) (*Company, error)
}

type ProfessionalRepository interface {
	CreateProfessional(ctx context.Context, professional *Professional) error
	GetProfessionalByID(ctx context.Context, id uuid.UUID) (*Professional, error)
	GetProfessionalByUserID(ctx context.Context, userID uuid.UUID) (*Professional, error)
	ListByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*Professional, error)
}
