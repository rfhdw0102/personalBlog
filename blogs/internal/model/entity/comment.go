package entity

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        int     `gorm:"primaryKey;autoIncrement"`
	ArticleID int     `gorm:"not null;index:idx_article_status_created"`
	Article   Article `gorm:"foreignKey:ArticleID;references:ID"`

	UserID int  `gorm:"index"`
	User   User `gorm:"foreignKey:UserID;references:ID"`

	Content   string         `gorm:"type:text;not null"`
	ParentID  int            `gorm:"index"`
	Status    int            `gorm:"type:tinyint;default:1;index:idx_article_status_created"` // 1显示 0隐藏
	CreatedAt time.Time      `gorm:"index:idx_article_status_created"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ParentComment *Comment `gorm:"foreignKey:ParentID;references:ID"`
}

func (comment *Comment) TableName() string {
	return "comments"
}
