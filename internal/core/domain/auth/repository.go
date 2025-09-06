package auth

import (
	"context"
	"youmeet/internal/core/domain/user"
)

// Service interface para autenticação
type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (*user.User, error)
	Login(ctx context.Context, req *AuthRequest) (*AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*user.User, error)
}