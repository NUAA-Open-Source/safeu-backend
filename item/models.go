package item

import (
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
	DownCount    string
	Type         string
	IsPublic     bool
	IsGroup      bool
}
