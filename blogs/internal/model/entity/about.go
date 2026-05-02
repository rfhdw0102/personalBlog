package entity

type About struct {
	Base
	Content   string `gorm:"type:longtext" json:"content"`
	Intro     string `gorm:"type:longtext" json:"intro"`
	TechStack string `gorm:"type:varchar(500)" json:"tech_stack"`
	Email     string `gorm:"type:varchar(100)" json:"email"`
	Github    string `gorm:"type:varchar(255)" json:"github"`
}

func (About) TableName() string {
	return "about"
}
