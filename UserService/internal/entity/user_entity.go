package entity

type User struct {
	Id       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:255;not null"`
	Email    string `gorm:"size:255;unique;not null"`
	Password string `gorm:"size:255;not null"`
	Token    string `gorm:"size:255"`
}
