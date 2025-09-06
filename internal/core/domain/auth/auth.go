package auth

import "youmeet/internal/core/domain/user"

// AuthRequest representa uma solicitação de autenticação
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse representa uma resposta de autenticação
type AuthResponse struct {
	Token string     `json:"token"`
	User  *user.User `json:"user"`
}

// RegisterRequest representa uma solicitação de registro
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}