package repositories

import (
	"context"
	"github.com/google/uuid"
	"youmeet/internal/core/domain/user"
)

type UserRepository struct {
	db DBClient
}

func NewUserRepository(db DBClient) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	return r.db.Create(u)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var u user.User
	err := r.db.First(&u, "id = ?", id)
	return &u, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.First(&u, "email = ?", email)
	return &u, err
}

type CompanyRepository struct {
	db DBClient
}

func NewCompanyRepository(db DBClient) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) CreateCompany(ctx context.Context, company *user.Company) error {
	return r.db.Create(company)
}

func (r *CompanyRepository) GetCompanyByID(ctx context.Context, id uuid.UUID) (*user.Company, error) {
	var company user.Company
	err := r.db.First(&company, "id = ?", id)
	return &company, err
}

func (r *CompanyRepository) GetCompanyByUserID(ctx context.Context, userID uuid.UUID) (*user.Company, error) {
	var company user.Company
	err := r.db.First(&company, "user_id = ?", userID)
	return &company, err
}

type ProfessionalRepository struct {
	db DBClient
}

func NewProfessionalRepository(db DBClient) *ProfessionalRepository {
	return &ProfessionalRepository{db: db}
}

func (r *ProfessionalRepository) CreateProfessional(ctx context.Context, professional *user.Professional) error {
	return r.db.Create(professional)
}

func (r *ProfessionalRepository) GetProfessionalByID(ctx context.Context, id uuid.UUID) (*user.Professional, error) {
	var professional user.Professional
	err := r.db.First(&professional, "id = ?", id)
	return &professional, err
}

func (r *ProfessionalRepository) GetProfessionalByUserID(ctx context.Context, userID uuid.UUID) (*user.Professional, error) {
	var professional user.Professional
	err := r.db.First(&professional, "user_id = ?", userID)
	return &professional, err
}

func (r *ProfessionalRepository) ListByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*user.Professional, error) {
	var professionals []*user.Professional
	err := r.db.Find(&professionals, "company_id = ?", companyID)
	return professionals, err
}