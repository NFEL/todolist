package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Email    string
	Password string
	Avatar   string
}

type UserSession struct {
	gorm.Model
	User   *User `gorm:"foreignKey:UserID"`
	UserID uint
	Valid  bool `gorm:"default:true;"`
}
