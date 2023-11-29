package models

import (
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name         string `gorm:"type:varchar(100);not null" json:"name"`
	Introduction string `gorm:"type:varchar(100);not null" json:"introduction"`
}
