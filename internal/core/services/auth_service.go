package services

import (
	"context"
	"errors"
	"time"

	"youmeet/internal/core/domain/auth"
	"youmeet/internal/core/domain/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    user.UserRepository
	companyRepo user.CompanyRepository
	profRepo    user.ProfessionalRepository
}

func NewAuthService(userRepo user.UserRepository, companyRepo user.CompanyRepository, profRepo user.ProfessionalRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		companyRepo: companyRepo,
		profRepo:    profRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, req *auth.RegisterRequest) (*user.User, error) {
	name, email, password, role := req.Name, req.Email, req.Password, req.Role
	// Verificar se usuário já existe
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("usuário já existe")
	}

	// Validar role
	if role != "client" && role != "company" && role != "professional" {
		return nil, errors.New("role inválido")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Criar usuário
	u := &user.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		CreatedAt:    time.Now(),
	}

	err = s.userRepo.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	// Criar perfis específicos baseado no role
	switch role {
	case "company":
		company := &user.Company{
			ID:        uuid.New(),
			UserID:    u.ID,
			Name:      name,
			CreatedAt: time.Now(),
		}
		s.companyRepo.CreateCompany(ctx, company)

	case "professional":
		professional := &user.Professional{
			ID:        uuid.New(),
			UserID:    u.ID,
			Name:      name,
			CompanyID: nil, // Autônomo
			CreatedAt: time.Now(),
		}
		s.profRepo.CreateProfessional(ctx, professional)
	}

	return u, nil
}

func (s *AuthService) Login(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {
	email, password := req.Email, req.Password
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Gerar token simples (em produção usar JWT)
	token := uuid.New().String()

	return &auth.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*user.User, error) {
	// Implementação simples - em produção validar JWT
	return nil, errors.New("token inválido")
}
