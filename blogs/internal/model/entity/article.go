package entity

import (
	"gorm.io/gorm"
)

type Article struct {
	Base
	Title      string `gorm:"type:varchar(255);not null"`
	Content    string `gorm:"type:longtext;not null"`
	Summary    string `gorm:"type:varchar(255)"`
	CoverImage string `gorm:"type:varchar(255)"`
	UserID     int    `gorm:"index:idx_user_status_updated"`
	User       User   `gorm:"foreignKey:UserID;references:ID"`

	CategoryID int      `gorm:"index:idx_status_category_created"`
	Category   Category `gorm:"foreignKey:CategoryID;references:ID"`

	Status    string         `gorm:"type:enum('draft', 'published');default:'draft';index:idx_status_category_created;index:idx_user_status_updated"`
	ViewCount int            `gorm:"default:0;index"`
	LikeCount int            `gorm:"default:0"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Tags []Tag `gorm:"many2many:article_tags;foreignKey:ID;joinForeignKey:ArticleID;references:ID;joinReferences:TagID"`
}

func (article *Article) TableName() string {
	return "articles"
}

// BeforeDelete 钩子
func (article *Article) BeforeDelete(tx *gorm.DB) (err error) {
	// 删除文章关联的点赞记录
	if err := tx.Where("article_id = ?", article.ID).Delete(&Like{}).Error; err != nil {
		return err
	}
	// 删除文章关联的评论记录
	if err := tx.Where("article_id = ?", article.ID).Delete(&Comment{}).Error; err != nil {
		return err
	}
	// 删除文章与标签的关联关系
	if err := tx.Table("article_tags").Where("article_id = ?", article.ID).Delete(nil).Error; err != nil {
		return err
	}
	return nil
}
