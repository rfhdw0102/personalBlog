package service

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"
	"blogs/internal/repository"
	appErrors "blogs/pkg/errors"
	"blogs/pkg/logger"
	"blogs/pkg/utils"
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo  repository.UserRepository
	redisRepo repository.RedisRepository
}

func NewUserService(userRepo repository.UserRepository, redisRepo repository.RedisRepository) UserService {
	return &userService{userRepo: userRepo, redisRepo: redisRepo}
}

// GetByID 根据 ID 获取用户
func (s *userService) GetByID(id int) (entity.User, error) {
	return s.userRepo.GetByID(id)
}

// GetAuthor 获得作者信息
func (s *userService) GetAuthor() (response.AuthorResponse, error) {
	return s.userRepo.GetAuthor()
}

func (s *userService) Update(id int, req request.UpdateUserRequest, key string) error {
	updates, err := s.buildUserUpdates(req, key)
	if err != nil {
		return err
	}
	if len(updates) == 0 {
		return nil
	}
	if req.Status == 0 {
		s.kickOffline(id)
	}

	if err := s.userRepo.Updates(id, updates); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "更新用户失败", err)
	}
	return nil
}

func (s *userService) Delete(id int) error {
	if err := s.userRepo.Delete(id); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "删除用户失败", err)
	}
	return nil
}

func (s *userService) List(req request.UserListRequest) ([]entity.User, int64, error) {
	return s.userRepo.List(req)
}

func (s *userService) UpdateAvatar(id int, avatarPath string) error {
	updates := make(map[string]interface{})
	updates["avatar"] = avatarPath
	if err := s.userRepo.Updates(id, updates); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "更新头像失败", err)
	}
	return nil
}

func (s *userService) processPassword(encryptedPwd, key string) (string, error) {
	password, err := utils.DecryptPassword(encryptedPwd, key)
	if err != nil {
		logger.Error("Password decryption failed", zap.Error(err))
		return "", appErrors.New(appErrors.CodeInternalError, "服务器内部错误")
	}

	if err := utils.ValidatePassword(password); err != nil {
		return "", appErrors.New(appErrors.CodeBadRequest, err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Bcrypt hashing failed", zap.Error(err))
		return "", appErrors.New(appErrors.CodeInternalError, "服务器内部错误")
	}
	return string(hash), nil
}

func (s *userService) Create(req request.CreateUserRequest, key string) error {
	if err := s.ensureUserUnique(req.Username, req.Email); err != nil {
		return err
	}

	if req.Password != req.ConfirmPassword {
		return appErrors.New(appErrors.CodeBadRequest, "密码不一致")
	}

	hashedPwd, err := s.processPassword(req.Password, key)
	if err != nil {
		return err
	}

	user := entity.User{
		Username: req.Username,
		Password: hashedPwd,
		Email:    req.Email,
		Role:     req.Role,
		Avatar:   "/uploads/avatars/default.png",
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return appErrors.Wrap(appErrors.CodeInternalError, "创建用户失败", err)
	}
	return nil
}

func (s *userService) ConsumePrivateKey(keyID string) (string, error) {
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

func (s *userService) GetStats() (response.AdminCard, error) {
	return s.userRepo.GetStats()
}

func (s *userService) buildUserUpdates(req request.UpdateUserRequest, key string) (map[string]interface{}, error) {
	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Password != "" {
		hashedPwd, err := s.processPassword(req.Password, key)
		if err != nil {
			return nil, err
		}
		updates["password"] = hashedPwd
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	updates["status"] = req.Status
	if req.Introduction != "" {
		updates["introduction"] = req.Introduction
	}
	if req.Qr != "" {
		updates["qr"] = req.Qr
	}
	return updates, nil
}

func (s *userService) kickOffline(id int) {
	keyRedis := fmt.Sprintf("user_token:%d", id)
	if err := s.redisRepo.DelKey(context.Background(), keyRedis); err != nil {
		logger.Warn("删除Token失败", zap.Int("id", id))
	}
}

func (s *userService) ensureUserUnique(username, email string) error {
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
