package entity

import (
	"gorm.io/gorm"
)

type NotificationType string

type Notification struct {
	Base
	RecipientID int  `gorm:"not null;index"`
	Recipient   User `gorm:"foreignKey:RecipientID;references:ID"`

	ActorID int  `gorm:"not null;index"`
	Actor   User `gorm:"foreignKey:ActorID;references:ID"`

	ArticleID int     `gorm:"not null;index"`
	Article   Article `gorm:"foreignKey:ArticleID;references:ID"`

	CommentID int      `gorm:"default:0"`
	Comment   *Comment `gorm:"foreignKey:CommentID;references:ID"`

	Type      NotificationType `gorm:"type:enum('like', 'comment', 'reply');not null"`
	IsRead    bool             `gorm:"default:false"`
	DeletedAt gorm.DeletedAt   `gorm:"index"`
}

func (Notification) TableName() string {
	return "notifications"
}
