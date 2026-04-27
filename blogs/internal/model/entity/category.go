package entity

import "gorm.io/gorm"

type Category struct {
	Base
	Name     string    `gorm:"type:varchar(50);not null" json:"name"`
	Articles []Article `gorm:"foreignKey:CategoryID;references:ID" json:"-"`
}

func (category *Category) TableName() string {
	return "category"
}

// BeforeDelete 钩子
func (category *Category) BeforeDelete(tx *gorm.DB) (err error) {
	// 将属于该分类的所有文章的 category_id 重置为 0
	return tx.Model(&Article{}).Where("category_id = ?", category.ID).Update("category_id", 0).Error
}
