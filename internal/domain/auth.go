package domain

import "context"

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type AuthUsecase interface {
	Register(ctx context.Context, req *RegisterRequest) error
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	GoogleLogin(ctx context.Context) (string, error)
	GoogleCallback(ctx context.Context, code string) (*AuthResponse, error)
	VerifyEmail(ctx context.Context, token string) error
}