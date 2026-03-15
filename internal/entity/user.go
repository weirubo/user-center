package entity

import "time"

type User struct {
	ID              int64     `gorm:"primaryKey" json:"id"`
	Email           string    `gorm:"uniqueIndex;size:255" json:"email"`
	Phone           string    `gorm:"index;size:20" json:"phone"`
	Password        string    `gorm:"size:255" json:"-"`
	Nickname        string    `gorm:"size:100" json:"nickname"`
	PasswordErrors   int       `gorm:"default:0" json:"password_errors"`
	LockedUntil     *time.Time `gorm:"index" json:"locked_until"`
}

func (User) TableName() string {
	return "users"
}
