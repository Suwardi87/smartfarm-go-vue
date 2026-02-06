package models

import (
	"time"
	"gorm.io/gorm"
)

type ProductView struct {
	gorm.Model
	ProductID uint
	UserID    uint // 0 if guest
	ViewedAt  time.Time
}
