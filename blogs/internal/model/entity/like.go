package entity

import (
	"time"
)

type Like struct {
	ID        int     `gorm:"primaryKey;autoIncrement"`
	ArticleID int     `gorm:"not null;uniqueIndex:idx_article_user"`
	Article   Article `gorm:"foreignKey:ArticleID;references:ID"`

	UserID int  `gorm:"not null;uniqueIndex:idx_article_user"`
	User   User `gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time
}

func (like *Like) TableName() string {
	return "likes"
}
