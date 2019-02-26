package item

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Item struct {
	gorm.Model
	Status       int    `gorm:"not null"`
	Name         string `gorm:"index"`
	OriginalName string `gorm:"not null"`
	Host         string
	ReCode       string `gorm:"index"`
	Password     string
	DownCount    int `gorm:"default:-100"`
	Type         string
	IsPublic     bool `gorm:"default:true"`
	ArchiveType  int  `gorm:"default:0"`
	Protocol     string
	Bucket       string
	Endpoint     string
	Path         string
	ExpiredAt    time.Time `gorm:"NOT NULL"`
}

type Token struct {
	gorm.Model
	Token        string `gorm:"NOT NULL;INDEX"`
	RetrieveCode string `gorm:"NOT NULL"`
	Valid        bool
	ExpiredAt    time.Time `gorm:"NOT NULL"`
}
