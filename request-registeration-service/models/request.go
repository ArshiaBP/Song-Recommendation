package models

import (
	"gorm.io/gorm"
)

type RequestInfo struct {
	gorm.Model
	Email  string `gorm:"unique;not null"`
	Status string `gorm:"not null"`
	SongID string `gorm:""`
}
