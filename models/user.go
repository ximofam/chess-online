package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Username  string    `json:"username" gorm:"type:varchar(30);uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null"`
	Role      string    `json:"role" gorm:"type:varchar(20);default:'user'"`
	CreatedAt time.Time `json:"createt_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
