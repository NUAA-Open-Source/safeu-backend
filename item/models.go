package item

import (
	"time"
)

type Item struct {
	ID           uint   `gorm:"AUTO_INCREMENT"`
	Status       int    `gorm:"not null"`
	Name         string `gorm:"primary_key"`
	OriginalName string `gorm:"not null"`
	Host         string
	ReCode       string `gorm:"index"`
	Password     string
	DownCount    string
	Type         string
	IsPublic     bool
	IsGroup      bool
	ExpiredTime  time.Time
	CreateTime   time.Time
	UpdatedTime  time.Time
}
