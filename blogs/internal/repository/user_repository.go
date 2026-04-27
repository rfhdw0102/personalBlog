package repository

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/dto/response"
	"blogs/internal/model/entity"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(user entity.User) error {
	return r.db.Create(&user).Error
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(email string) (entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

// UpdatePassword 更新密码
func (r *userRepository) UpdatePassword(email, password string) error {
	return r.db.Model(&entity.User{}).Where("email = ?", email).Update("password", password).Error
}

// GetByID 根据 ID 获取用户
func (r *userRepository) GetByID(id int) (entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return user, err
}

// GetAuthor 获取作者信息
func (r *userRepository) GetAuthor() (response.AuthorResponse, error) {
	var author response.AuthorResponse
	err := r.db.Table("users").
		Select(`
			users.id,
			users.username,
			users.avatar,
			users.introduction,
			users.qr,
			COALESCE(COUNT(articles.id), 0) as article_count,
			COALESCE(SUM(articles.view_count), 0) as view_count
		`).
		Joins("LEFT JOIN articles ON users.id = articles.user_id AND articles.deleted_at IS NULL").
		Where("users.id = ?", 1).
		Group("users.id").
		First(&author).Error
	return author, err
}

// Updates 局部更新用户信息
func (r *userRepository) Updates(id int, values map[string]interface{}) error {
	return r.db.Model(&entity.User{}).Where("id = ?", id).Updates(values).Error
}

// Delete 删除用户
func (r *userRepository) Delete(id int) error {
	// 使用模型实例进行删除，以触发 GORM 的 BeforeDelete 钩子实现级联删除
	return r.db.Delete(&entity.User{Base: entity.Base{ID: id}}).Error
}

// List 分页获取用户列表
func (r *userRepository) List(req request.UserListRequest) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	db := r.db.Model(&entity.User{})
	if req.Query != "" {
		db = db.Where("username LIKE ? OR email LIKE ?", "%"+req.Query+"%", "%"+req.Query+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Select("id, username, email, avatar, role, status, introduction, created_at, updated_at").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&users).Error
	return users, total, err
}

func (r *userRepository) GetStats() (response.AdminCard, error) {
	var userCount int64
	var articleCount int64
	var categoryCount int64
	var tagCount int64
	if err := r.db.Model(&entity.User{}).Count(&userCount).Error; err != nil {
		return response.AdminCard{}, err
	}
	if err := r.db.Model(&entity.Article{}).Count(&articleCount).Error; err != nil {
		return response.AdminCard{}, err
	}
	if err := r.db.Model(&entity.Category{}).Count(&categoryCount).Error; err != nil {
		return response.AdminCard{}, err
	}
	if err := r.db.Model(&entity.Tag{}).Count(&tagCount).Error; err != nil {
		return response.AdminCard{}, err
	}
	return response.AdminCard{UserCount: int(userCount), ArticleCount: int(articleCount), CategoryCount: int(categoryCount), TagCount: int(tagCount)}, nil
}
