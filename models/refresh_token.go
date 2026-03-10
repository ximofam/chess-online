package models

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Revoked   bool      `gorm:"not null;default:false"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
}
