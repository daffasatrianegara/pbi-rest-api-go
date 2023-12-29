package models

import (
	"time"
)

type User struct {
	ID        uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    Username  string         `json:"username" gorm:"unique;not null"`
    Email     string         `json:"email" gorm:"unique;not null"`
    Password  string         `json:"password" gorm:"not null"`
    Photos    []Photo        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}