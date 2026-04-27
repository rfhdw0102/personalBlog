package entity

import "gorm.io/gorm"

type Tag struct {
	Base
	Name string `gorm:"size:50;unique;not null" json:"name"`
}

func (tag *Tag) TableName() string {
	return "tags"
}

// BeforeDelete 钩子
func (tag *Tag) BeforeDelete(tx *gorm.DB) (err error) {
	// 删除 article_tags 中该标签的所有关联记录
	return tx.Table("article_tags").Where("tag_id = ?", tag.ID).Delete(nil).Error
}
