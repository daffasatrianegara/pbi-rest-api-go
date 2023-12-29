package models

import (
    "time"
)

type Photo struct {
    ID        uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    Title     string         `gorm:"column:title;not null" json:"title"`
    Caption   string         `gorm:"column:caption;not null" json:"caption"`
    PhotoUrl  string         `gorm:"column:photo_url;not null" json:"photo_url"`
    UserID    uint           `gorm:"column:user_id;not null" json:"user_id"`
    User      User           `gorm:"foreignKey:UserID" json:"user"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}