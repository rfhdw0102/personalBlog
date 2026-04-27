package service

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"
	"blogs/internal/repository"
	"blogs/pkg/config"
	appErrors "blogs/pkg/errors"
	"blogs/pkg/jwt"
	"blogs/pkg/logger"
	"blogs/pkg/utils"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo  repository.UserRepository
	redisRepo repository.RedisRepository
}

func NewAuthService(userRepo repository.UserRepository, redisRepo repository.RedisRepository) AuthService {
	return &authService{userRepo: userRepo, redisRepo: redisRepo}
}

func (s *authService) Register(req request.AuthRequest, key string) error {
	if err := s.validateVerifyCode(req.Email, req.Code); err != nil {
		return err
	}
	if err := s.ensureUserUnique(req.Username, req.Email); err != nil {
		return err
	}

	password, err := s.resolvePassword(req.Password, req.ConfirmPassword, key)
	if err != nil {
		return err
	}

	user := entity.User{
		Username: req.Username,
		Password: password,
		Email:    req.Email,
		Avatar:   "/uploads/avatars/default.png",
	}
	if err := s.userRepo.Create(user); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "注册失败", err)
	}

	s.clearVerifyCode(req.Email)
	return nil
}

func (s *authService) Login(req request.LoginRequest, key string) (response.LoginUserResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return response.LoginUserResponse{}, appErrors.New(appErrors.CodeNotFound, "用户不存在")
	}
	if user.Status == 0 {
		return response.LoginUserResponse{}, appErrors.New(appErrors.CodeForbidden, "账户已被禁用")
	}

	password, err := utils.DecryptPassword(req.Password, key)
	if err != nil {
		return response.LoginUserResponse{}, appErrors.New(appErrors.CodeInternalError, "服务器内部错误")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return response.LoginUserResponse{}, appErrors.New(appErrors.CodeBadRequest, "密码错误")
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return response.LoginUserResponse{}, appErrors.New(appErrors.CodeInternalError, "生成 token 失败")
	}

	ctx := context.Background()
	keyRedis := fmt.Sprintf("user_token:%d", user.ID)
	expireTime := config.Get().JWT.ExpireHours * time.Hour
	if err = s.redisRepo.SetKey(ctx, keyRedis, token, expireTime); err != nil {
		logger.Error("存入redis token失败", zap.Error(err))
		return response.LoginUserResponse{}, appErrors.New(appErrors.CodeInternalError, "服务器内部错误")
	}

	return s.buildLoginResponse(user, token), nil
}

func (s *authService) Logout(ctx context.Context, userID int) error {
	key := fmt.Sprintf("user_token:%d", userID)
	return s.redisRepo.DelKey(ctx, key)
}

func (s *authService) SendCode(email string) error {
	code, err := s.generateCode()
	if err != nil {
		return appErrors.New(appErrors.CodeInternalError, "生成验证码失败")
	}

	ctx := context.Background()
	if err = s.redisRepo.SetKey(ctx, "code:"+email, code, 5*time.Minute); err != nil {
		return appErrors.New(appErrors.CodeInternalError, "发送验证码失败")
	}

	subject := "验证码"
	body := fmt.Sprintf("您的验证码为：%s，有效期为 5 分钟。", code)
	if err = utils.SendEmail(email, subject, body); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "发送验证码失败", err)
	}
	return nil
}

func (s *authService) ResetPassword(req request.ForgotPasswordRequest, key string) error {
	if err := s.validateVerifyCode(req.Email, req.Code); err != nil {
		return err
	}
	if err := s.ensureEmailRegistered(req.Email); err != nil {
		return err
	}

	password, err := s.resolvePassword(req.Password, req.ConfirmPassword, key)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(req.Email, password); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "重置密码失败", err)
	}

	s.clearVerifyCode(req.Email)
	return nil
}

func (s *authService) ConsumePrivateKey(keyID string) (string, error) {
	key, err := s.redisRepo.GetKey(context.Background(), "rsa:"+keyID)
	if err != nil {
		logger.Error("获取私钥失败", zap.Error(err))
		return "", appErrors.New(appErrors.CodeInternalError, "获取私钥失败")
	}

	if err := s.redisRepo.DelKey(context.Background(), "rsa:"+keyID); err != nil {
		logger.Warn("删除私钥失败", zap.Error(err))
		return "", appErrors.New(appErrors.CodeInternalError, "获取私钥失败")
	}

	return key, nil
}

func (s *authService) GenerateRSAKey() (response.RsaKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return response.RsaKey{}, appErrors.New(appErrors.CodeInternalError, "密钥生成失败")
	}

	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateBytes})
	keyID := uuid.New().String()
	if err = s.redisRepo.SetKey(context.Background(), "rsa:"+keyID, string(privatePEM), 5*time.Minute); err != nil {
		logger.Error("私钥存入redis失败", zap.Error(err))
		return response.RsaKey{}, appErrors.New(appErrors.CodeInternalError, "服务器内部错误")
	}

	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return response.RsaKey{}, appErrors.New(appErrors.CodeInternalError, "公钥序列化失败")
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubASN1})
	return response.RsaKey{KeyID: keyID, PubKey: string(pubPEM)}, nil
}

func (s *authService) buildLoginResponse(user entity.User, token string) response.LoginUserResponse {
	return response.LoginUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Role:      user.Role,
		Status:    user.Status,
		Token:     token,
	}
}

func (s *authService) validateVerifyCode(email, inputCode string) error {
	code, err := s.redisRepo.GetKey(context.Background(), "code:"+email)
	if err != nil {
		return appErrors.New(appErrors.CodeBadRequest, "验证码已过期或不存在")
	}
	if code != inputCode {
		return appErrors.New(appErrors.CodeBadRequest, "验证码错误")
	}
	return nil
}

func (s *authService) generateCode() (string, error) {
	codeNum, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", codeNum.Int64()+100000), nil
}

func (s *authService) ensureUserUnique(username, email string) error {
	exists, err := s.userRepo.ExistsByUsername(username)
	if err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "查询用户失败", err)
	}
	if exists {
		return appErrors.New(appErrors.CodeBadRequest, "用户已存在")
	}

	exists, err = s.userRepo.ExistsByEmail(email)
	if err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "查询邮箱失败", err)
	}
	if exists {
		return appErrors.New(appErrors.CodeBadRequest, "邮箱已被注册")
	}
	return nil
}

func (s *authService) ensureEmailRegistered(email string) error {
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "查询邮箱失败", err)
	}
	if !exists {
		return appErrors.New(appErrors.CodeNotFound, "邮箱未注册")
	}
	return nil
}

func (s *authService) resolvePassword(passwordText, confirmText, key string) (string, error) {
	password, err := utils.DecryptPassword(passwordText, key)
	if err != nil {
		return "", appErrors.New(appErrors.CodeInternalError, "服务器内部错误")
	}
	if err := utils.ValidatePassword(password); err != nil {
		return "", appErrors.New(appErrors.CodeBadRequest, err.Error())
	}

	confirmPassword, err := utils.DecryptPassword(confirmText, key)
	if err != nil {
		return "", appErrors.New(appErrors.CodeInternalError, "服务器内部错误")
	}
	if password != confirmPassword {
		return "", appErrors.New(appErrors.CodeBadRequest, "密码不一致")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Bcrypt hashing failed", zap.Error(err))
		return "", appErrors.New(appErrors.CodeInternalError, "服务器内部错误")
	}
	return string(hash), nil
}

func (s *authService) clearVerifyCode(email string) {
	if err := s.redisRepo.DelKey(context.Background(), "code:"+email); err != nil {
		logger.Warn("删除缓存验证码失败", zap.String("email", email))
	}
}
