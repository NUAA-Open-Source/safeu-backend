package event

import "github.com/jinzhu/gorm"

type Event struct {
	gorm.Model
	Name string `gorm:"NOT NULL"`
	From string
}
