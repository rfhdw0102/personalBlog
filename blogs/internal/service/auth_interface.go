package service

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/dto/response"
	"context"
)

type AuthService interface {
	Register(req request.AuthRequest, key string) error
	Login(req request.LoginRequest, key string) (response.LoginUserResponse, error)
	Logout(ctx context.Context, userID int) error
	SendCode(email string) error
	ResetPassword(req request.ForgotPasswordRequest, key string) error
	ConsumePrivateKey(keyID string) (string, error)
	GenerateRSAKey() (response.RsaKey, error)
}
