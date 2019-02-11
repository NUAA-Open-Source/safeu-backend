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
	DownCount    int `gorm:"default:-1"`
	Type         string
	IsPublic     bool
	IsGroup      bool
}

type Token struct {
	gorm.Model
	Token        string `gorm:"NOT NULL;INDEX"`
	RetrieveCode string `gorm:"NOT NULL"`
	Valid        bool
	ExpiredAt    time.Time `gorm:"NOT NULL"`
}
