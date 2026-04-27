package entity

import (
	"gorm.io/gorm"
)

type User struct {
	Base
	Username     string    `gorm:"size:20;not null;unique" json:"username"`
	Password     string    `gorm:"size:200;not null" json:"-"`
	Email        string    `gorm:"size:50" json:"email"`
	Avatar       string    `gorm:"size:255" json:"avatar"`
	Role         string    `gorm:"type:enum('admin','user');default:'user'" json:"role"`
	Status       int       `gorm:"default:1" json:"status"`
	Introduction string    `gorm:"type:text" json:"introduction"`
	Qr           string    `gorm:"size:255" json:"qr"`
	Articles     []Article `gorm:"foreignKey:UserID;references:ID"`
}

func (user *User) TableName() string {
	return "users"
}

// BeforeDelete 钩子
func (user *User) BeforeDelete(tx *gorm.DB) (err error) {
	// 获取该用户的所有文章
	var articles []Article
	if err := tx.Where("user_id = ?", user.ID).Find(&articles).Error; err != nil {
		return err
	}
	// 逐个删除文章，以触发 Article 的 BeforeDelete 钩子
	for _, article := range articles {
		// 这里通过模型实例进行删除，以确保 Article 关联的点赞、评论和标签关联也被级联删除
		if err := tx.Delete(&article).Error; err != nil {
			return err
		}
	}
	return nil
}
